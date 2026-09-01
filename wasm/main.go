//go:build wasm

package main

import (
	"encoding/json"
	"syscall/js"
	"wasmcut/shared"
)

// Global timeline
var timeline *shared.Timeline

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

	clipID := "clip-1" // Simplified for v1

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

func main() {
	// Register functions to be callable from JavaScript
	js.Global().Set("createProject", js.FuncOf(createProject))
	js.Global().Set("getTimelineState", js.FuncOf(getTimelineState))
	js.Global().Set("addClip", js.FuncOf(addClip))
	js.Global().Set("undo", js.FuncOf(undo))
	js.Global().Set("redo", js.FuncOf(redo))

	// Keep the Go runtime alive so JavaScript can call the registered functions
	select {}
}
