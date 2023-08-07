package producer

import "github.com/evermos/boilerplate-go/event/model"

// Producer represents an event producer interface.
type Producer interface {
	Publish(request model.PublishRequest) error
}
