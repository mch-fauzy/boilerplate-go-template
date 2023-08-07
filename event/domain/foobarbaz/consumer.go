package foobarbaz

import (
	"encoding/json"
	"net/http"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/consumer"
	"github.com/evermos/boilerplate-go/event/model"
	"github.com/evermos/boilerplate-go/internal/domain/foobarbaz"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/rs/zerolog/log"
)

// ConsumerImpl is the SQS consumer implementation for this domain.
type ConsumerImpl struct {
	Config   *configs.Config
	Service  foobarbaz.FooService
	Consumer consumer.Consumer
}

// ProvideConsumerImpl is the provider for this consumer.
func ProvideConsumerImpl(config *configs.Config, service foobarbaz.FooService) ConsumerImpl {
	c := ConsumerImpl{}
	c.Config = config
	c.Service = service

	sqsConsumer := consumer.NewSQSConsumer(config)
	sqsConsumer.Process = c.processEvent
	c.Consumer = sqsConsumer

	return c
}

// Start starts up the SQS subscriber
func (c *ConsumerImpl) Start() {
	if c.Config.Event.Consumer.SQS.Topics.FooBarBaz.Enabled {
		go c.Consumer.Listen(c.Config.Event.Consumer.SQS.Topics.FooBarBaz.URL)
	}
}

func (c *ConsumerImpl) processEvent(value []byte) (err error) {
	snsMessage := model.SNSMessage{}
	err = json.Unmarshal(value, &snsMessage)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	log.
		Info().
		Str("topicARN", snsMessage.TopicARN).
		Interface("value", snsMessage).
		Msg("Received SNS message")

	requestFormat := foobarbaz.FooRequestFormat{}
	err = json.Unmarshal([]byte(snsMessage.Message), &requestFormat)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	_, err = c.Service.Create(requestFormat, snsMessage.MessageID)
	if err != nil {
		err = c.checkError(err)
	}

	return
}

func (c *ConsumerImpl) checkError(err error) error {
	f, ok := err.(*failure.Failure)
	if ok {
		if f.Code == http.StatusBadRequest {
			err = nil
		}
	}

	if err != nil {
		logger.ErrorWithStack(err)
	}

	return err
}
