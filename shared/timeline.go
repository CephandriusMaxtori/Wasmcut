package shared

// Project represents the root editing project
type Project struct {
	ID         string      `json:"id"`
	Version    int         `json:"version"`
	Tracks     []Track     `json:"tracks"`
	Transitions []Transition `json:"transitions"`
	Effects    []Effect    `json:"effects"`
}

// Track represents a video or audio track
type Track struct {
	ID    string `json:"id"`
	Type  string `json:"type"` // "video" or "audio"
	Clips []Clip `json:"clips"`
}

// Clip represents a segment of media on a track
type Clip struct {
	ID       string  `json:"id"`
	MediaRef string  `json:"media_ref"` // media://source.mp4
	In       float64 `json:"in"`        // trim start in source
	Out      float64 `json:"out"`       // trim end in source
	Position float64 `json:"position"` // timeline position
}

// Transition represents a transition between clips
type Transition struct {
	ID       string  `json:"id"`
	ClipA    string  `json:"clip_a"`
	ClipB    string  `json:"clip_b"`
	Type     string  `json:"type"`     // "cut", "crossfade", "wipe"
	Duration float64 `json:"duration"` // in seconds
}

// Effect represents an effect/filter on a clip or track
type Effect struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	TargetID   string            `json:"target_id"` // clip or track ID
	Parameters map[string]interface{} `json:"parameters"`
}
