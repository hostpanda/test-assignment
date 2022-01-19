package ws

import "encoding/json"

type Event string

const (
	SubscribeEvent   Event = "subscribe"
	UnsubscribeEvent Event = "unsubscribe"
)

// Payload is a data struct for a message.
type Payload struct {
	Event    Event           `json:"event"`
	Channels []string        `json:"channels"`
	Data     json.RawMessage `json:"data"`
}
