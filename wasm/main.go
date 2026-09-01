//go:build wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"wasmcut/shared"
)

// Global timeline
var timeline *shared.Timeline
var nextClipID int

func init() {
	// Initialize timeline on load
	timeline = shared.NewTimeline()
	timeline.AddTrack("track-1", "video")
}

// ResponseWrapper wraps results for the host
type ResponseWrapper struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// createProject creates a new timeline project
func createProject(this js.Value, args []js.Value) interface{} {
	timeline = shared.NewTimeline()
	timeline.AddTrack("track-1", "video")
	nextClipID = 1

	response := ResponseWrapper{
		Success: true,
		Data: map[string]interface{}{
			"project_id": timeline.Project.ID,
			"tracks":     len(timeline.Project.Tracks),
		},
	}

	bytes, _ := json.Marshal(response)
	return string(bytes)
}

// getTimelineState returns the current project state as JSON
func getTimelineState(this js.Value, args []js.Value) interface{} {
	if timeline == nil || timeline.Project == nil {
		response := ResponseWrapper{
			Success: false,
			Error:   "No project loaded",
		}
		bytes, _ := json.Marshal(response)
		return string(bytes)
	}

	response := ResponseWrapper{
		Success: true,
		Data:    timeline.Project,
	}

	bytes, _ := json.Marshal(response)
	return string(bytes)
}

// addClip adds a clip to a track
func addClip(this js.Value, args []js.Value) interface{} {
	if timeline == nil || len(args) < 5 {
		response := ResponseWrapper{
			Success: false,
			Error:   "Invalid arguments or no project",
		}
		bytes, _ := json.Marshal(response)
		return string(bytes)
	}

	trackID := args[0].String()
	mediaRef := args[1].String()
	inTime := args[2].Float()
	outTime := args[3].Float()
	position := args[4].Float()

	clipID := fmt.Sprintf("clip-%d", nextClipID)
	nextClipID++

	cmd := &shared.AddClipCmd{
		ID:      clipID,
		TrackID: trackID,
		Clip: shared.Clip{
			ID:       clipID,
			MediaRef: mediaRef,
			In:       inTime,
			Out:      outTime,
			Position: position,
		},
	}

	err := timeline.ExecuteCommand(cmd)
	response := ResponseWrapper{
		Success: err == nil,
		Error:   errorStr(err),
		Data: map[string]interface{}{
			"clip_id": clipID,
		},
	}

	bytes, _ := json.Marshal(response)
	return string(bytes)
}

// trimClip updates a clip's source in/out points.
func trimClip(this js.Value, args []js.Value) interface{} {
	if timeline == nil || len(args) < 6 {
		return jsonResponse(false, "Invalid arguments or no project", nil)
	}

	cmd := &shared.TrimClipCmd{
		ClipID: args[0].String(), TrackID: args[1].String(),
		OldIn: args[2].Float(), OldOut: args[3].Float(),
		NewIn: args[4].Float(), NewOut: args[5].Float(),
	}
	err := timeline.ExecuteCommand(cmd)
	return jsonResponse(err == nil, errorStr(err), nil)
}

// deleteClip removes a clip from the timeline.
func deleteClip(this js.Value, args []js.Value) interface{} {
	if timeline == nil || len(args) < 2 {
		return jsonResponse(false, "Invalid arguments or no project", nil)
	}
	cmd := &shared.DeleteClipCmd{ClipID: args[0].String(), TrackID: args[1].String()}
	err := timeline.ExecuteCommand(cmd)
	return jsonResponse(err == nil, errorStr(err), nil)
}

// moveClip changes a clip's timeline position or track.
func moveClip(this js.Value, args []js.Value) interface{} {
	if timeline == nil || len(args) < 5 {
		return jsonResponse(false, "Invalid arguments or no project", nil)
	}
	cmd := &shared.MoveClipCmd{
		ClipID: args[0].String(), OldTrackID: args[1].String(), NewTrackID: args[2].String(),
		OldPosition: args[3].Float(), NewPosition: args[4].Float(),
	}
	err := timeline.ExecuteCommand(cmd)
	return jsonResponse(err == nil, errorStr(err), nil)
}

// undo undoes the last action
func undo(this js.Value, args []js.Value) interface{} {
	if timeline == nil {
		response := ResponseWrapper{
			Success: false,
			Error:   "No project loaded",
		}
		bytes, _ := json.Marshal(response)
		return string(bytes)
	}

	err := timeline.Undo()
	response := ResponseWrapper{
		Success: err == nil,
		Error:   errorStr(err),
	}

	bytes, _ := json.Marshal(response)
	return string(bytes)
}

// redo redoes the last undone action
func redo(this js.Value, args []js.Value) interface{} {
	if timeline == nil {
		response := ResponseWrapper{
			Success: false,
			Error:   "No project loaded",
		}
		bytes, _ := json.Marshal(response)
		return string(bytes)
	}

	err := timeline.Redo()
	response := ResponseWrapper{
		Success: err == nil,
		Error:   errorStr(err),
	}

	bytes, _ := json.Marshal(response)
	return string(bytes)
}

// Helper functions
func errorStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func jsonResponse(success bool, message string, data interface{}) string {
	response := ResponseWrapper{Success: success, Error: message, Data: data}
	bytes, _ := json.Marshal(response)
	return string(bytes)
}

func main() {
	// Register functions to be callable from JavaScript
	js.Global().Set("createProject", js.FuncOf(createProject))
	js.Global().Set("getTimelineState", js.FuncOf(getTimelineState))
	js.Global().Set("addClip", js.FuncOf(addClip))
	js.Global().Set("trimClip", js.FuncOf(trimClip))
	js.Global().Set("deleteClip", js.FuncOf(deleteClip))
	js.Global().Set("moveClip", js.FuncOf(moveClip))
	js.Global().Set("undo", js.FuncOf(undo))
	js.Global().Set("redo", js.FuncOf(redo))

	// Keep the Go runtime alive so JavaScript can call the registered functions
	select {}
}
