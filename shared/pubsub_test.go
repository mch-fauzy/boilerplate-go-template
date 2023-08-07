package shared_test

import (
	"errors"
	"testing"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expected := "Testing"
		var actual string
		pubsub := shared.New(1, shared.SetMessageBuffer(10))
		pubsub.SubscriberRegistry("test", func(message []byte) error {
			actual = string(message)
			return nil
		})
		pubsub.Start()
		pubsub.Publish("test", []byte("Testing"))

		time.Sleep(1 * time.Second)
		assert.Equal(t, expected, actual)
	})

	t.Run("Failed Message", func(t *testing.T) {
		expected := "Testing"
		var actual string
		pubsub := shared.New(1, shared.SetMessageBuffer(10))
		pubsub.SubscriberRegistry("test", func(message []byte) error {
			actual = string(message)
			return nil
		})
		pubsub.Start()
		pubsub.Publish("test", []byte("b"))

		time.Sleep(1 * time.Second)
		assert.NotEqual(t, expected, actual)
	})

	t.Run("Retry Success", func(t *testing.T) {
		expectedTopicTest := 5
		expectedTopicTest2 := 2
		retryCountTest := 0
		retryCountTest2 := 0
		pubsub := shared.New(1, shared.SetMessageBuffer(10))
		pubsub.SubscriberRegistry("test", func(message []byte) error {
			retryCountTest++
			return errors.New("error test retry")
		}, shared.SetMaxRetry(5))
		pubsub.SubscriberRegistry("test-2", func(message []byte) error {
			retryCountTest2++
			return errors.New("error test retry")
		}, shared.SetMaxRetry(2))

		pubsub.Start()
		pubsub.Publish("test", []byte("b"))
		pubsub.Publish("test-2", []byte("b"))

		time.Sleep(1 * time.Second)
		assert.Equal(t, expectedTopicTest, retryCountTest)
		assert.Equal(t, expectedTopicTest2, retryCountTest2)
	})

	t.Run("Test Race", func(t *testing.T) {
		counter := 0
		pubsub := shared.New(10)
		pubsub.SubscriberRegistry("test", func(message []byte) error {
			counter++
			return nil
		}, shared.SetAsynchronousThread(true))
		pubsub.Start()

		for i := 0; i < 1000; i++ {
			pubsub.Publish("test", []byte("test"))
		}

		time.Sleep(3 * time.Second)
		assert.Equal(t, 1000, counter)
	})
}
