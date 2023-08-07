package producer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/event/model"
	"github.com/rs/zerolog/log"
)

// SNSConfig is a wrapper for SNS configuration.
type SNSConfig struct {
	Config *configs.Config
}

func createCredentials(config *configs.Config) credentials.Value {
	return credentials.Value{
		AccessKeyID:     config.Event.Producer.SNS.AccessKeyID,
		SecretAccessKey: config.Event.Producer.SNS.SecretAccessKey,
	}
}

func createSNSConfig(config *configs.Config) (sessionSNS *session.Session, err error) {
	sessionSNS, err = session.NewSession(&aws.Config{
		Region:      aws.String(config.Event.Producer.SNS.Region),
		Credentials: credentials.NewStaticCredentialsFromCreds(createCredentials(config)),
	})

	if err != nil {
		log.Fatal().Err(err).Msg("failed creating new SNS session")
		return
	}

	return
}

// SNSProducer is an SNS producer.
type SNSProducer struct {
	config *configs.Config
	sns    *sns.SNS
}

// NewSNSProducer creates a new object from Producer
func NewSNSProducer(config *configs.Config) *SNSProducer {
	sess, err := createSNSConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating SNS config")
	}
	log.Info().Str("region", *sess.Config.Region).Msg("SNS Producer ready to publish messages.")
	return &SNSProducer{config: config, sns: sns.New(sess)}
}

// Publish publishes a message to SNS.
func (p *SNSProducer) Publish(request model.PublishRequest) error {
	err := p.sendMessage(&sns.PublishInput{
		Message:        aws.String(string(request.Event.Data.Value)),
		MessageGroupId: request.MessageGroupID,
		TopicArn:       &request.Topic,
	})

	return err
}

func (p *SNSProducer) sendMessage(msg *sns.PublishInput) error {
	resp, err := p.sns.Publish(msg)
	if err != nil {
		log.Err(err).Interface("output", *msg).Msg("failed publishing message")
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
