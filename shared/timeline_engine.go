package shared

import "errors"

// CommandStack manages undo/redo operations
type CommandStack struct {
	commands []Command
	index    int // points to the next command to redo
}

// NewCommandStack creates a new empty command stack
func NewCommandStack() *CommandStack {
	return &CommandStack{
		commands: []Command{},
		index:    0,
	}
}

// Execute runs a command and adds it to the stack
func (cs *CommandStack) Execute(cmd Command, project *Project) error {
	if err := cmd.Execute(project); err != nil {
		return err
	}

	// Truncate any commands after current index (remove redo history)
	cs.commands = cs.commands[:cs.index]

	// Add new command
	cs.commands = append(cs.commands, cmd)
	cs.index = len(cs.commands)

	return nil
}

// Undo reverts the last command
func (cs *CommandStack) Undo(project *Project) error {
	if cs.index == 0 {
		return ErrEmptyStack
	}

	cs.index--
	cmd := cs.commands[cs.index]
	return cmd.Undo(project)
}

// Redo re-applies the last undone command
func (cs *CommandStack) Redo(project *Project) error {
	if cs.index >= len(cs.commands) {
		return ErrEmptyStack
	}

	cmd := cs.commands[cs.index]
	if err := cmd.Execute(project); err != nil {
		return err
	}

	cs.index++
	return nil
}

// CanUndo returns true if there are commands to undo
func (cs *CommandStack) CanUndo() bool {
	return cs.index > 0
}

// CanRedo returns true if there are commands to redo
func (cs *CommandStack) CanRedo() bool {
	return cs.index < len(cs.commands)
}

// GetCommandCount returns the total number of commands in history
func (cs *CommandStack) GetCommandCount() int {
	return len(cs.commands)
}

// GetCurrentIndex returns the current position in command history
func (cs *CommandStack) GetCurrentIndex() int {
	return cs.index
}

// Timeline wraps a Project with command handling
type Timeline struct {
	Project *Project
	Stack   *CommandStack
}

// NewTimeline creates a new timeline with an empty project
func NewTimeline() *Timeline {
	return &Timeline{
		Project: &Project{
			ID:          "project-1",
			Version:     1,
			Tracks:      []Track{},
			Transitions: []Transition{},
			Effects:     []Effect{},
		},
		Stack: NewCommandStack(),
	}
}

// AddTrack adds a new track to the timeline
func (t *Timeline) AddTrack(id string, trackType string) error {
	for _, track := range t.Project.Tracks {
		if track.ID == id {
			return errors.New("track with this ID already exists")
		}
	}

	t.Project.Tracks = append(t.Project.Tracks, Track{
		ID:    id,
		Type:  trackType,
		Clips: []Clip{},
	})
	return nil
}

// ExecuteCommand executes a command through the command stack
func (t *Timeline) ExecuteCommand(cmd Command) error {
	return t.Stack.Execute(cmd, t.Project)
}

// Undo undoes the last command
func (t *Timeline) Undo() error {
	return t.Stack.Undo(t.Project)
}

// Redo redoes the last undone command
func (t *Timeline) Redo() error {
	return t.Stack.Redo(t.Project)
}

// GetUndoStack returns the command stack
func (t *Timeline) GetUndoStack() *CommandStack {
	return t.Stack
}
