//go:build wasm

package main

import (
	"encoding/json"
	"wasmcut/shared"
)

// Global timeline
var timeline *shared.Timeline
var lastError error

func init() {
	// Initialize timeline on load
	timeline = shared.NewTimeline()
	timeline.AddTrack("track-1", "video")
}

// Helper: serialize and store result
func storeResult(data interface{}) []byte {
	bytes, err := json.Marshal(data)
	if err != nil {
		lastError = err
		return nil
	}
	return bytes
}

// ResponseWrapper wraps results for the host
type ResponseWrapper struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

//export CreateProject
func CreateProject() int32 {
	timeline = shared.NewTimeline()
	timeline.AddTrack("track-1", "video")

	response := ResponseWrapper{
		Success: true,
		Data: map[string]interface{}{
			"project_id": timeline.Project.ID,
			"tracks":     len(timeline.Project.Tracks),
		},
	}
	data := storeResult(response)
	return int32(len(data))
}

//export GetTimelineState
func GetTimelineState() int32 {
	if timeline == nil || timeline.Project == nil {
		return 0
	}

	response := ResponseWrapper{
		Success: true,
		Data:    timeline.Project,
	}
	data := storeResult(response)
	return int32(len(data))
}

//export AddClip
func AddClip(trackID string, mediaRef string, inTime float64, outTime float64, position float64) int32 {
	if timeline == nil {
		return -1
	}

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
	data := storeResult(response)
	return int32(len(data))
}

//export Undo
func Undo() int32 {
	if timeline == nil {
		return -1
	}

	err := timeline.Undo()
	response := ResponseWrapper{
		Success: err == nil,
		Error:   errorStr(err),
	}
	data := storeResult(response)
	return int32(len(data))
}

//export Redo
func Redo() int32 {
	if timeline == nil {
		return -1
	}

	err := timeline.Redo()
	response := ResponseWrapper{
		Success: err == nil,
		Error:   errorStr(err),
	}
	data := storeResult(response)
	return int32(len(data))
}

// Helper functions
func errorStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func getTrackByID(id string) *shared.Track {
	if timeline == nil || timeline.Project == nil {
		return nil
	}
	for i := range timeline.Project.Tracks {
		if timeline.Project.Tracks[i].ID == id {
			return &timeline.Project.Tracks[i]
		}
	}
	return nil
}

func main() {
	// Keep the Go runtime alive so JavaScript can call the exported functions.
	select {}
}
