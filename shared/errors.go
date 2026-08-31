package shared

import "errors"

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrTrackNotFound   = errors.New("track not found")
	ErrClipNotFound    = errors.New("clip not found")
	ErrEmptyStack      = errors.New("stack is empty")
)
