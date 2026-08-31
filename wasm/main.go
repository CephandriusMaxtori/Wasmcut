//go:build wasm

package main

import (
	"encoding/json"
	"fmt"
	"wasmcut/shared"
)

// Global project store (simplified for v1)
var currentProject *shared.Project

//export CreateProject
func CreateProject() *uint8 {
	currentProject = &shared.Project{
		ID:      "project-1",
		Version: 1,
		Tracks: []shared.Track{
			{
				ID:    "track-1",
				Type:  "video",
				Clips: []shared.Clip{},
			},
		},
		Transitions: []shared.Transition{},
		Effects:     []shared.Effect{},
	}
	result, err := json.Marshal(currentProject)
	if err != nil {
		fmt.Printf("Error marshaling project: %v\n", err)
		return nil
	}
	// Return pointer to the first byte of the JSON
	return &result[0]
}

//export GetTimelineState
func GetTimelineState() *uint8 {
	if currentProject == nil {
		return nil
	}
	result, err := json.Marshal(currentProject)
	if err != nil {
		fmt.Printf("Error marshaling timeline: %v\n", err)
		return nil
	}
	return &result[0]
}

//export AddClip
func AddClip(trackID string, mediaRef string, in float64, out float64, position float64) bool {
	if currentProject == nil {
		return false
	}
	for i, track := range currentProject.Tracks {
		if track.ID == trackID {
			clip := shared.Clip{
				ID:       fmt.Sprintf("clip-%d", len(track.Clips)+1),
				MediaRef: mediaRef,
				In:       in,
				Out:      out,
				Position: position,
			}
			currentProject.Tracks[i].Clips = append(currentProject.Tracks[i].Clips, clip)
			return true
		}
	}
	return false
}

func main() {
	// WASI entry point - required but not used in this architecture
	// Wasm module is called via host ABI only
}
