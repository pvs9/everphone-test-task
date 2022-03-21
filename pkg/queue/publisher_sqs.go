package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
)

type PublisherSQS struct {
	client *Client
}

func NewPublisherSQS(client *Client) *PublisherSQS {
	return &PublisherSQS{client: client}
}

func (p *PublisherSQS) Publish(messageBody string) (*string, error) {
	result, err := p.client.Client.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    p.client.QueueUrl,
		MessageBody: aws.String(messageBody),
	})

	if err != nil {
		return nil, err
	}

	messageId := result.MessageId
	log.Infof("message with ID: %s successfully published", *messageId)

	return messageId, nil
}
