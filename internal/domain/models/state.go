package models

// State represents a specific state in the conversation flow
type State string

// StateGroup represents a group of related states
type StateGroup string

// Group returns the group this state belongs to
func (s State) Group() StateGroup {
	// Extract group from state name (format: "group:state")
	// If no group specified, return empty group
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			return StateGroup(s[:i])
		}
	}
	return ""
}

// DefaultState is used when no state is set
const DefaultState State = ""
