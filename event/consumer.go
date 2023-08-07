package event

import (
	"github.com/evermos/boilerplate-go/event/domain/foobarbaz"
)

// Consumers is the wrapper to contain all event consumers.
type Consumers struct {
	FooBarBaz foobarbaz.ConsumerImpl
}

// ProvideConsumers is the provider function for Consumers.
func ProvideConsumers(fooBarBaz foobarbaz.ConsumerImpl) Consumers {
	return Consumers{
		FooBarBaz: fooBarBaz,
	}
}

// Start starts all domains event consumer
func (c *Consumers) Start() {
	c.FooBarBaz.Start()
}
