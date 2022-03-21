package queue

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Consumer interface {
	Consume(handler func(m *sqs.Message) error)
	DefaultHandler(m *sqs.Message) error
}

type Publisher interface {
	Publish(messageBody string) (*string, error)
}

type Client struct {
	Client         *sqs.SQS
	ConsumerConfig *ConsumerConfig
	QueueUrl       *string
}

type ConsumerConfig struct {
	MaxNumberOfMessages      int64
	MessageVisibilityTimeout int64
	PollDelayInMilliseconds  int
	Receivers                int
}

type Queue struct {
	Consumer
	Publisher
}

func NewQueue(client *Client) *Queue {
	return &Queue{
		Consumer:  NewConsumerSQS(client),
		Publisher: NewPublisherSQS(client),
	}
}
