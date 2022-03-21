package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type ReceiverSQS struct {
	client          *Client
	messagesChannel chan []*sqs.Message
	shutdown        chan os.Signal
}

func (r *ReceiverSQS) applyBackPressure() {
	time.Sleep(time.Millisecond * time.Duration(r.client.ConsumerConfig.PollDelayInMilliseconds))
}

func (r *ReceiverSQS) receiveMessages() {
	for {

		select {
		case <-r.shutdown:
			log.Print("Shutting down message receiver")
			close(r.messagesChannel)
			return
		default:
			result, err := r.client.Client.ReceiveMessage(&sqs.ReceiveMessageInput{
				AttributeNames: []*string{
					aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
				},
				MessageAttributeNames: []*string{
					aws.String(sqs.QueueAttributeNameAll),
				},
				QueueUrl:            r.client.QueueUrl,
				MaxNumberOfMessages: aws.Int64(r.client.ConsumerConfig.MaxNumberOfMessages),
				VisibilityTimeout:   aws.Int64(r.client.ConsumerConfig.MessageVisibilityTimeout),
			})

			if err != nil {
				log.Errorf("could not read from queue: %s", err)
				return
			}

			if len(result.Messages) > 0 {
				messages := result.Messages
				r.messagesChannel <- messages
			}

			r.applyBackPressure()
		}
	}
}
