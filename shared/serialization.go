package shared

import "encoding/json"

// SaveProject serializes the project to JSON bytes
func (t *Timeline) SaveProject() ([]byte, error) {
	return json.MarshalIndent(t.Project, "", "  ")
}

// LoadProject deserializes a project from JSON bytes
func (t *Timeline) LoadProject(data []byte) error {
	project := &Project{}
	if err := json.Unmarshal(data, project); err != nil {
		return err
	}
	t.Project = project
	// Reset command stack on load
	t.Stack = NewCommandStack()
	return nil
}

// LoadProjectJSON loads a project from a JSON string
func (t *Timeline) LoadProjectJSON(jsonStr string) error {
	return t.LoadProject([]byte(jsonStr))
}

// GetProjectJSON returns the project as a JSON string
func (t *Timeline) GetProjectJSON() (string, error) {
	data, err := t.SaveProject()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ProjectSnapshot represents the state at a specific point in history
type ProjectSnapshot struct {
	Project       Project `json:"project"`
	CommandIndex  int     `json:"command_index"`
	TotalCommands int     `json:"total_commands"`
}

// GetSnapshot returns a snapshot of the current project state
func (t *Timeline) GetSnapshot() ProjectSnapshot {
	return ProjectSnapshot{
		Project:       *t.Project,
		CommandIndex:  t.Stack.GetCurrentIndex(),
		TotalCommands: t.Stack.GetCommandCount(),
	}
}
