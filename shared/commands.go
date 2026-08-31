package shared

// Command defines the interface for all edit operations
type Command interface {
	// Execute performs the operation
	Execute(project *Project) error
	// Undo reverses the operation
	Undo(project *Project) error
	// GetDescription returns a human-readable name
	GetDescription() string
}

// AddClipCmd adds a new clip to a track
type AddClipCmd struct {
	ID      string
	TrackID string
	Clip    Clip
}

func (c *AddClipCmd) Execute(p *Project) error {
	for i, track := range p.Tracks {
		if track.ID == c.TrackID {
			// Store original state for undo
			c.Clip.ID = c.ID
			p.Tracks[i].Clips = append(p.Tracks[i].Clips, c.Clip)
			return nil
		}
	}
	return ErrTrackNotFound
}

func (c *AddClipCmd) Undo(p *Project) error {
	for i, track := range p.Tracks {
		if track.ID == c.TrackID {
			// Remove the clip we added
			for j, clip := range track.Clips {
				if clip.ID == c.ID {
					p.Tracks[i].Clips = append(track.Clips[:j], track.Clips[j+1:]...)
					return nil
				}
			}
		}
	}
	return ErrClipNotFound
}

func (c *AddClipCmd) GetDescription() string {
	return "Add Clip"
}

// TrimClipCmd modifies the in/out points of a clip
type TrimClipCmd struct {
	ClipID  string
	TrackID string
	OldIn   float64
	OldOut  float64
	NewIn   float64
	NewOut  float64
}

func (c *TrimClipCmd) Execute(p *Project) error {
	for _, track := range p.Tracks {
		if track.ID == c.TrackID {
			for i, clip := range track.Clips {
				if clip.ID == c.ClipID {
					track.Clips[i].In = c.NewIn
					track.Clips[i].Out = c.NewOut
					return nil
				}
			}
		}
	}
	return ErrClipNotFound
}

func (c *TrimClipCmd) Undo(p *Project) error {
	for _, track := range p.Tracks {
		if track.ID == c.TrackID {
			for i, clip := range track.Clips {
				if clip.ID == c.ClipID {
					track.Clips[i].In = c.OldIn
					track.Clips[i].Out = c.OldOut
					return nil
				}
			}
		}
	}
	return ErrClipNotFound
}

func (c *TrimClipCmd) GetDescription() string {
	return "Trim Clip"
}

// MoveClipCmd moves a clip to a new position and/or track
type MoveClipCmd struct {
	ClipID      string
	OldTrackID  string
	NewTrackID  string
	OldPosition float64
	NewPosition float64
}

func (c *MoveClipCmd) Execute(p *Project) error {
	var clip *Clip
	var oldTrackIdx int

	// Find and remove from old track
	for i, track := range p.Tracks {
		if track.ID == c.OldTrackID {
			for j, cl := range track.Clips {
				if cl.ID == c.ClipID {
					clip = &p.Tracks[i].Clips[j]
					oldTrackIdx = i
					break
				}
			}
			if clip != nil {
				break
			}
		}
	}

	if clip == nil {
		return ErrClipNotFound
	}

	// Update position
	clip.Position = c.NewPosition

	// If moving to different track, remove and re-add
	if c.OldTrackID != c.NewTrackID {
		for i, track := range p.Tracks {
			if track.ID == c.NewTrackID {
				tempClip := *clip
				p.Tracks[oldTrackIdx].Clips = append(
					p.Tracks[oldTrackIdx].Clips[:indexOfClip(p.Tracks[oldTrackIdx].Clips, c.ClipID)],
					p.Tracks[oldTrackIdx].Clips[indexOfClip(p.Tracks[oldTrackIdx].Clips, c.ClipID)+1:]...,
				)
				p.Tracks[i].Clips = append(p.Tracks[i].Clips, tempClip)
				return nil
			}
		}
		return ErrTrackNotFound
	}

	return nil
}

func (c *MoveClipCmd) Undo(p *Project) error {
	var clip *Clip
	var newTrackIdx int

	// Find clip in new location
	for i, track := range p.Tracks {
		if track.ID == c.NewTrackID {
			for j, cl := range track.Clips {
				if cl.ID == c.ClipID {
					clip = &p.Tracks[i].Clips[j]
					newTrackIdx = i
					break
				}
			}
			if clip != nil {
				break
			}
		}
	}

	if clip == nil {
		return ErrClipNotFound
	}

	// Restore old position
	clip.Position = c.OldPosition

	// If was moved between tracks, reverse it
	if c.OldTrackID != c.NewTrackID {
		for i, track := range p.Tracks {
			if track.ID == c.OldTrackID {
				tempClip := *clip
				p.Tracks[newTrackIdx].Clips = append(
					p.Tracks[newTrackIdx].Clips[:indexOfClip(p.Tracks[newTrackIdx].Clips, c.ClipID)],
					p.Tracks[newTrackIdx].Clips[indexOfClip(p.Tracks[newTrackIdx].Clips, c.ClipID)+1:]...,
				)
				p.Tracks[i].Clips = append(p.Tracks[i].Clips, tempClip)
				return nil
			}
		}
		return ErrTrackNotFound
	}

	return nil
}

func (c *MoveClipCmd) GetDescription() string {
	return "Move Clip"
}

// DeleteClipCmd removes a clip from the timeline
type DeleteClipCmd struct {
	ClipID  string
	TrackID string
	// Store deleted clip for undo
	DeletedClip Clip
	DeletedIdx  int
}

func (c *DeleteClipCmd) Execute(p *Project) error {
	for i, track := range p.Tracks {
		if track.ID == c.TrackID {
			for j, clip := range track.Clips {
				if clip.ID == c.ClipID {
					c.DeletedClip = clip
					c.DeletedIdx = j
					p.Tracks[i].Clips = append(track.Clips[:j], track.Clips[j+1:]...)
					return nil
				}
			}
		}
	}
	return ErrClipNotFound
}

func (c *DeleteClipCmd) Undo(p *Project) error {
	for i, track := range p.Tracks {
		if track.ID == c.TrackID {
			// Restore at original index if possible
			if c.DeletedIdx <= len(track.Clips) {
				newClips := make([]Clip, len(track.Clips)+1)
				copy(newClips[:c.DeletedIdx], track.Clips[:c.DeletedIdx])
				newClips[c.DeletedIdx] = c.DeletedClip
				copy(newClips[c.DeletedIdx+1:], track.Clips[c.DeletedIdx:])
				p.Tracks[i].Clips = newClips
			} else {
				p.Tracks[i].Clips = append(track.Clips, c.DeletedClip)
			}
			return nil
		}
	}
	return ErrTrackNotFound
}

func (c *DeleteClipCmd) GetDescription() string {
	return "Delete Clip"
}

// Helper function
func indexOfClip(clips []Clip, clipID string) int {
	for i, clip := range clips {
		if clip.ID == clipID {
			return i
		}
	}
	return -1
}
