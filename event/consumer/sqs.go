package consumer

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/rs/zerolog/log"
)

// Process represents the processing function of the message consumer.
type Process func(e []byte) error

// SQSConfig represents an SQS configuration object.
type SQSConfig struct {
	Config configs.Config
}

// IsExpired checks whether the configuration is expired.
func (m *SQSConfig) IsExpired() bool {
	return false
}

// Retrieve performs an SQS retrieve function.
func (m *SQSConfig) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     m.Config.Event.Consumer.SQS.AccessKeyID,
		SecretAccessKey: m.Config.Event.Consumer.SQS.SecretAccessKey,
	}, nil
}

func createSQSConfig(config *configs.Config) (*session.Session, error) {
	sqsConfig := SQSConfig{Config: *config}
	return session.NewSession(&aws.Config{
		Region:      &config.Event.Consumer.SQS.Region,
		Credentials: credentials.NewCredentials(&sqsConfig),
		MaxRetries:  aws.Int(config.Event.Consumer.SQS.MaxRetries),
	})
}

// SQSConsumer represents an SQS consumer.
type SQSConsumer struct {
	Process Process
	config  *configs.Config
	sqs     *sqs.SQS
}

// NewSQSConsumer create object Consumer
func NewSQSConsumer(config *configs.Config) *SQSConsumer {
	sess, err := createSQSConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating sqs config")
	}
	return &SQSConsumer{config: config, sqs: sqs.New(sess)}
}

// Listen is a function to listen new message from sqs queue
func (p *SQSConsumer) Listen(url string) {
	log.Info().Str("url", url).Msg("SQS Consumer will start polling.")

	retries := 0
	for {
		receiveResp, err := p.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(url),
			MaxNumberOfMessages: aws.Int64(p.config.Event.Consumer.SQS.MaxMessage),
			WaitTimeSeconds:     aws.Int64(p.config.Event.Consumer.SQS.WaitTimeSeconds),
		})
		if err != nil {
			if retries == p.config.Event.Consumer.SQS.MaxRetriesConsume {
				log.Error().Err(err).Int("retries", retries).Msg("failed receiving message after maximum retries, failing permanently")
				return
			}

			log.
				Error().
				Err(err).
				Str("url", url).
				Int("retries", retries).
				Int("backoffSeconds", p.config.Event.Consumer.SQS.BackoffSeconds).
				Msg("failed receiving message, will retry")
			retries++
			time.Sleep(time.Duration(p.config.Event.Consumer.SQS.BackoffSeconds) * time.Second)
			continue
		} else {
			retries = 0
		}

		for _, message := range receiveResp.Messages {
			err := p.Process([]byte(*message.Body))
			if err != nil {
				log.Error().Err(err).Msg("failed processing message")
			}

			err = p.deleteMessage(message, url)
			if err != nil {
				log.Error().Err(err).Msg("failed deleting message")
			}
		}
	}
}

func (p *SQSConsumer) deleteMessage(msg *sqs.Message, url string) error {
	output, err := p.sqs.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &url,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Err(err).Interface("output", output).Msg("failed deleting message")
		return err
	}
	return nil
}
