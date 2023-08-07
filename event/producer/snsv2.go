package producer

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/model"

	"github.com/rs/zerolog/log"
)

type SNSProducerV2 struct {
	cfg    *configs.Config
	client *sns.Client
}

func NewSnsProducerV2(cfg *configs.Config, httpClient *http.Client) SNSProducerV2 {
	credentials := credentials.NewStaticCredentialsProvider(cfg.Event.Producer.SNS.AccessKeyID, cfg.Event.Consumer.SQS.SecretAccessKey, "")
	config, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(cfg.Event.Producer.SNS.Region),
		config.WithCredentialsProvider(credentials),
		config.WithRetryer(func() aws.Retryer {
			return retry.AddWithMaxAttempts(retry.NewStandard(), cfg.Event.Producer.SNS.MaxRetries)
		}),
	)

	if httpClient != nil {
		config.HTTPClient = httpClient
	}

	if err != nil {
		log.Fatal().Msg("error instantiate sns producer")
	}

	return SNSProducerV2{
		client: sns.NewFromConfig(config),
		cfg:    cfg,
	}
}

func (s *SNSProducerV2) Publish(request model.PublishRequest) error {
	return s.publish(request)
}

func (s *SNSProducerV2) publish(request model.PublishRequest) error {
	msg := &sns.PublishInput{
		Message:        aws.String(string(request.Event.Data.Value)),
		MessageGroupId: request.MessageGroupID,
		TopicArn:       &request.Topic,
	}

	resp, err := s.client.Publish(context.TODO(), msg)
	if err != nil {
		return err
	}

	logMsg := log.Info()

	if msg.PhoneNumber != nil {
		logMsg.Str("phoneNumber", *msg.PhoneNumber)
	}

	if msg.TargetArn != nil {
		logMsg.Str("targetArn", *msg.TargetArn)
	}

	if msg.TopicArn != nil {
		logMsg.Str("topicArn", *msg.TopicArn)
	}

	if resp != nil && resp.MessageId != nil {
		logMsg.Str("messageId", *resp.MessageId)
	}

	logMsg.
		Interface("snsMessage", msg.Message).
		Msg("Published SNS message")

	return nil
}
