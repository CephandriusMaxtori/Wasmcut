package shared

import (
	"encoding/json"
	"testing"
)

func TestTimelineCreate(t *testing.T) {
	tl := NewTimeline()
	if tl.Project == nil {
		t.Fatal("Expected project to be created")
	}
	if tl.Project.ID != "project-1" {
		t.Errorf("Expected project ID 'project-1', got %s", tl.Project.ID)
	}
	if tl.Project.Version != 1 {
		t.Errorf("Expected version 1, got %d", tl.Project.Version)
	}
}

func TestAddClip(t *testing.T) {
	tl := NewTimeline()
	tl.AddTrack("track-1", "video")

	cmd := &AddClipCmd{
		ID:      "clip-1",
		TrackID: "track-1",
		Clip: Clip{
			ID:       "clip-1",
			MediaRef: "media://test.mp4",
			In:       0.0,
			Out:      5.0,
			Position: 0.0,
		},
	}

	err := tl.ExecuteCommand(cmd)
	if err != nil {
		t.Fatalf("ExecuteCommand failed: %v", err)
	}

	track := tl.Project.Tracks[0]
	if len(track.Clips) != 1 {
		t.Errorf("Expected 1 clip, got %d", len(track.Clips))
	}

	clip := track.Clips[0]
	if clip.ID != "clip-1" {
		t.Errorf("Expected clip ID 'clip-1', got %s", clip.ID)
	}
	if clip.MediaRef != "media://test.mp4" {
		t.Errorf("Expected media ref 'media://test.mp4', got %s", clip.MediaRef)
	}
}

func TestTrimClip(t *testing.T) {
	tl := NewTimeline()
	tl.AddTrack("track-1", "video")

	// Add a clip
	addCmd := &AddClipCmd{
		ID:      "clip-1",
		TrackID: "track-1",
		Clip: Clip{
			ID:       "clip-1",
			MediaRef: "media://test.mp4",
			In:       0.0,
			Out:      10.0,
			Position: 0.0,
		},
	}
	tl.ExecuteCommand(addCmd)

	// Trim the clip
	trimCmd := &TrimClipCmd{
		ClipID:  "clip-1",
		TrackID: "track-1",
		OldIn:   0.0,
		OldOut:  10.0,
		NewIn:   2.0,
		NewOut:  8.0,
	}
	err := tl.ExecuteCommand(trimCmd)
	if err != nil {
		t.Fatalf("TrimClip failed: %v", err)
	}

	clip := tl.Project.Tracks[0].Clips[0]
	if clip.In != 2.0 || clip.Out != 8.0 {
		t.Errorf("Expected In:2.0, Out:8.0, got In:%f, Out:%f", clip.In, clip.Out)
	}
}

func TestMoveClip(t *testing.T) {
	tl := NewTimeline()
	tl.AddTrack("track-1", "video")
	tl.AddTrack("track-2", "video")

	// Add clip to track-1
	addCmd := &AddClipCmd{
		ID:      "clip-1",
		TrackID: "track-1",
		Clip: Clip{
			ID:       "clip-1",
			MediaRef: "media://test.mp4",
			In:       0.0,
			Out:      5.0,
			Position: 0.0,
		},
	}
	tl.ExecuteCommand(addCmd)

	// Move to track-2
	moveCmd := &MoveClipCmd{
		ClipID:      "clip-1",
		OldTrackID:  "track-1",
		NewTrackID:  "track-2",
		OldPosition: 0.0,
		NewPosition: 10.0,
	}
	err := tl.ExecuteCommand(moveCmd)
	if err != nil {
		t.Fatalf("MoveClip failed: %v", err)
	}

	if len(tl.Project.Tracks[0].Clips) != 0 {
		t.Errorf("Expected 0 clips in track-1, got %d", len(tl.Project.Tracks[0].Clips))
	}

	if len(tl.Project.Tracks[1].Clips) != 1 {
		t.Errorf("Expected 1 clip in track-2, got %d", len(tl.Project.Tracks[1].Clips))
	}

	clip := tl.Project.Tracks[1].Clips[0]
	if clip.Position != 10.0 {
		t.Errorf("Expected position 10.0, got %f", clip.Position)
	}
}

func TestDeleteClip(t *testing.T) {
	tl := NewTimeline()
	tl.AddTrack("track-1", "video")

	// Add a clip
	addCmd := &AddClipCmd{
		ID:      "clip-1",
		TrackID: "track-1",
		Clip: Clip{
			ID:       "clip-1",
			MediaRef: "media://test.mp4",
			In:       0.0,
			Out:      5.0,
			Position: 0.0,
		},
	}
	tl.ExecuteCommand(addCmd)

	if len(tl.Project.Tracks[0].Clips) != 1 {
		t.Fatal("Expected 1 clip after add")
	}

	// Delete the clip
	delCmd := &DeleteClipCmd{
		ClipID:  "clip-1",
		TrackID: "track-1",
	}
	err := tl.ExecuteCommand(delCmd)
	if err != nil {
		t.Fatalf("DeleteClip failed: %v", err)
	}

	if len(tl.Project.Tracks[0].Clips) != 0 {
		t.Errorf("Expected 0 clips after delete, got %d", len(tl.Project.Tracks[0].Clips))
	}
}

func TestUndo(t *testing.T) {
	tl := NewTimeline()
	tl.AddTrack("track-1", "video")

	// Add a clip
	addCmd := &AddClipCmd{
		ID:      "clip-1",
		TrackID: "track-1",
		Clip: Clip{
			ID:       "clip-1",
			MediaRef: "media://test.mp4",
			In:       0.0,
			Out:      5.0,
			Position: 0.0,
		},
	}
	tl.ExecuteCommand(addCmd)
	if len(tl.Project.Tracks[0].Clips) != 1 {
		t.Fatal("Expected 1 clip after add")
	}

	// Undo
	err := tl.Undo()
	if err != nil {
		t.Fatalf("Undo failed: %v", err)
	}

	if len(tl.Project.Tracks[0].Clips) != 0 {
		t.Errorf("Expected 0 clips after undo, got %d", len(tl.Project.Tracks[0].Clips))
	}
}

func TestRedo(t *testing.T) {
	tl := NewTimeline()
	tl.AddTrack("track-1", "video")

	// Add a clip
	addCmd := &AddClipCmd{
		ID:      "clip-1",
		TrackID: "track-1",
		Clip: Clip{
			ID:       "clip-1",
			MediaRef: "media://test.mp4",
			In:       0.0,
			Out:      5.0,
			Position: 0.0,
		},
	}
	tl.ExecuteCommand(addCmd)

	// Undo
	tl.Undo()
	if len(tl.Project.Tracks[0].Clips) != 0 {
		t.Fatal("Expected 0 clips after undo")
	}

	// Redo
	err := tl.Redo()
	if err != nil {
		t.Fatalf("Redo failed: %v", err)
	}

	if len(tl.Project.Tracks[0].Clips) != 1 {
		t.Errorf("Expected 1 clip after redo, got %d", len(tl.Project.Tracks[0].Clips))
	}
}

func TestSerializationRoundtrip(t *testing.T) {
	// Create timeline with some clips
	tl1 := NewTimeline()
	tl1.AddTrack("track-1", "video")

	addCmd := &AddClipCmd{
		ID:      "clip-1",
		TrackID: "track-1",
		Clip: Clip{
			ID:       "clip-1",
			MediaRef: "media://test.mp4",
			In:       1.5,
			Out:      6.5,
			Position: 0.0,
		},
	}
	tl1.ExecuteCommand(addCmd)

	// Serialize
	data, err := tl1.SaveProject()
	if err != nil {
		t.Fatalf("SaveProject failed: %v", err)
	}

	// Deserialize
	tl2 := NewTimeline()
	err = tl2.LoadProject(data)
	if err != nil {
		t.Fatalf("LoadProject failed: %v", err)
	}

	// Verify
	if len(tl2.Project.Tracks) != 1 {
		t.Errorf("Expected 1 track, got %d", len(tl2.Project.Tracks))
	}

	if len(tl2.Project.Tracks[0].Clips) != 1 {
		t.Errorf("Expected 1 clip, got %d", len(tl2.Project.Tracks[0].Clips))
	}

	clip := tl2.Project.Tracks[0].Clips[0]
	if clip.ID != "clip-1" || clip.In != 1.5 || clip.Out != 6.5 {
		t.Error("Clip data mismatch after roundtrip")
	}
}

func TestSnapshotInfo(t *testing.T) {
	tl := NewTimeline()
	tl.AddTrack("track-1", "video")

	snapshot := tl.GetSnapshot()
	if snapshot.CommandIndex != 0 {
		t.Errorf("Expected command index 0, got %d", snapshot.CommandIndex)
	}

	if snapshot.TotalCommands != 0 {
		t.Errorf("Expected 0 total commands, got %d", snapshot.TotalCommands)
	}

	// Add a clip
	addCmd := &AddClipCmd{
		ID:      "clip-1",
		TrackID: "track-1",
		Clip: Clip{
			ID:       "clip-1",
			MediaRef: "media://test.mp4",
			In:       0.0,
			Out:      5.0,
			Position: 0.0,
		},
	}
	tl.ExecuteCommand(addCmd)

	snapshot = tl.GetSnapshot()
	if snapshot.CommandIndex != 1 {
		t.Errorf("Expected command index 1, got %d", snapshot.CommandIndex)
	}

	if snapshot.TotalCommands != 1 {
		t.Errorf("Expected 1 total command, got %d", snapshot.TotalCommands)
	}
}

func TestJSONOutput(t *testing.T) {
	tl := NewTimeline()
	tl.AddTrack("track-1", "video")

	jsonStr, err := tl.GetProjectJSON()
	if err != nil {
		t.Fatalf("GetProjectJSON failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}

	if result["id"] != "project-1" {
		t.Errorf("Expected project ID 'project-1' in JSON")
	}
}
