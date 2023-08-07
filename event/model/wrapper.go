package model

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

// SNSMessage is a wrapper struct for messages received in SQS that originated
// from SNS.
type SNSMessage struct {
	Type             string    `json:"Type"`
	MessageID        uuid.UUID `json:"MessageId"`
	TopicARN         string    `json:"TopicArn"`
	Message          string    `json:"Message"`
	Timestamp        string    `json:"Timestamp"`
	SignatureVersion string    `json:"SignatureVersion"`
	Signature        string    `json:"Signature"`
	SigningCertURL   string    `json:"SigningCertURL"`
	UnsubscribeURL   string    `json:"UnsubscribeURL"`
}

// EventWrapper is the wrapper object for events.
type EventWrapper struct {
	EventType string `json:"event_type"`
	Data      Data   `json:"data"`
}

// Data contains the data that is to be sent using an event.
type Data struct {
	Timestamp time.Time `json:"timestamp"`
	Value     []byte    `json:"value"`
}

// NewEvent creates a new event given an event type and an arbitrary model.
// Returns an EventWrapper object.
func NewEvent(eventType string, model interface{}) EventWrapper {
	value, _ := json.Marshal(model)

	return EventWrapper{
		EventType: eventType,
		Data: Data{
			Timestamp: time.Now(),
			Value:     value,
		},
	}
}

// PublishRequest is a wrapper for all message publishing requests.
type PublishRequest struct {
	Channel        string
	Event          EventWrapper
	MessageGroupID *string
	Topic          string
}
