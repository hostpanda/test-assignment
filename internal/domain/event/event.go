package event

import (
	"encoding/json"
)

// Event is a data struct for a message.
type Event struct {
	UserID  string          `json:"user_id,omitempty"`
	Command string          `json:"event,omitempty"`
	Channel string          `json:"channel,omitempty"`
	Message json.RawMessage `json:"message,omitempty"`
}
