package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Config struct {
	Options   session.Options
	QueueName string
	ConsumerConfig
}

func NewSQSQueue(cfg Config) (*Client, error) {
	sess, err := session.NewSessionWithOptions(cfg.Options)

	if err != nil {
		return nil, err
	}

	sqsClient := sqs.New(sess)

	queueUrlOutput, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(cfg.QueueName),
	})

	if err != nil {
		return nil, err
	}

	return &Client{sqsClient, &cfg.ConsumerConfig, queueUrlOutput.QueueUrl}, nil
}
