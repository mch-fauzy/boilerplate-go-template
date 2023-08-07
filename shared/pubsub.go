package shared

import (
	"time"
)

type message struct {
	topic   string
	payload []byte
}

type Consumer struct {
	message     chan message
	messagePool chan chan message
	runner      map[string]TopicRunner
}

type Process func(message []byte) error

type consumerConfig struct {
	MaxRetry           int
	MaxDelayRetry      time.Duration
	AsynchronousThread bool
}

func SetMaxRetry(maxRetry int) func(*consumerConfig) {
	return func(cc *consumerConfig) {
		cc.MaxRetry = maxRetry
	}
}

func SetMaxDelayRetry(maxDelay time.Duration) func(*consumerConfig) {
	return func(cc *consumerConfig) {
		cc.MaxDelayRetry = maxDelay
	}
}

// SetAsynchronousThread if enabled process will synchronously process by total flight
func SetAsynchronousThread(sync bool) func(*consumerConfig) {
	return func(cc *consumerConfig) {
		cc.AsynchronousThread = sync
	}
}

func defaultConsumerConfig() consumerConfig {
	return consumerConfig{
		MaxRetry:           0,
		MaxDelayRetry:      0 * time.Second,
		AsynchronousThread: false,
	}
}

func consumer(messagePol chan chan message, subscribers map[string]TopicRunner) Consumer {
	return Consumer{
		message:     make(chan message),
		messagePool: messagePol,
		runner:      subscribers,
	}
}

func (c Consumer) consume() {
	go func() {
		for {
			// starts out empty
			// send the response to dispatcher
			c.messagePool <- c.message

			// read the response
			msg := <-c.message

			runner := c.runner[msg.topic]
			if runner.consumerConfig.AsynchronousThread {
				go runner.backoff(func() error {
					return runner.Process(msg.payload)
				})
				continue
			}

			// if enabled process concurrent will handle by max flight
			runner.backoff(func() error {
				return runner.Process(msg.payload)
			})
		}
	}()
}

type PubSub struct {
	message     chan message
	messagePool chan chan message
	max         int
	topics      map[string]TopicRunner
}

type pubsubConfig struct {
	MessageBuffer int
}

func defaultPubsubConfig() pubsubConfig {
	return pubsubConfig{
		MessageBuffer: 0,
	}
}

func SetMessageBuffer(maxBuffer int) func(*pubsubConfig) {
	return func(pc *pubsubConfig) {
		pc.MessageBuffer = maxBuffer
	}
}

type TopicRunner struct {
	Process        Process
	consumerConfig consumerConfig
}

func (r TopicRunner) backoff(exec func() error) error {
	var err error

	if r.consumerConfig.MaxRetry == 0 {
		return exec()
	}

	counter := 0
	for counter < r.consumerConfig.MaxRetry {
		err = exec()
		if err == nil {
			return err
		}

		counter++
		time.Sleep(r.consumerConfig.MaxDelayRetry)
	}

	return err
}

// max is a total process could be handle
func New(maxFlight int, opts ...func(*pubsubConfig)) PubSub {
	config := defaultPubsubConfig()
	for _, opt := range opts {
		opt(&config)
	}

	return PubSub{
		message:     make(chan message, config.MessageBuffer),
		messagePool: make(chan chan message),
		max:         maxFlight,
		topics:      make(map[string]TopicRunner),
	}
}

func (p PubSub) Publish(topic string, payload []byte) {
	p.message <- message{
		topic:   topic,
		payload: payload,
	}
}

func (p PubSub) SubscriberRegistry(topicListener string, pr Process, opts ...func(*consumerConfig)) {
	cfg := defaultConsumerConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	p.topics[topicListener] = TopicRunner{
		Process:        pr,
		consumerConfig: cfg,
	}
}

func (p PubSub) Start() {
	for i := 0; i < p.max; i++ {
		consumer := consumer(p.messagePool, p.topics)
		consumer.consume()
	}

	go p.dispatch()
}

func (p PubSub) dispatch() {
	for {
		// waiting from p.Message from instantiate
		msg := <-p.message

		// read the response from consume
		response := <-p.messagePool

		// send value from function Publish to response / consumer
		response <- msg
	}
}
