package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"sync"
)

type ProcessorSQS struct {
	client  *Client
	handler func(*sqs.Message) error
}

func (p *ProcessorSQS) processMessages(messages []*sqs.Message) {
	nMessages := len(messages)
	deleteChannel := make(chan *string, nMessages)
	wg := sync.WaitGroup{}
	wg.Add(nMessages)

	for _, m := range messages {
		go func(message *sqs.Message) {
			defer wg.Done()
			err := p.handler(message)

			if err != nil {
				log.Errorf("error while handling message: %s", err)
				return
			}

			deleteChannel <- message.ReceiptHandle
		}(m)
	}

	wg.Wait()

	close(deleteChannel)
	entries := make([]*sqs.DeleteMessageBatchRequestEntry, 0, nMessages)

	for receipt := range deleteChannel {
		entries = append(entries, &sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(uuid.NewString()),
			ReceiptHandle: receipt,
		})
	}

	if len(entries) > 0 {
		_, dErr := p.client.Client.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
			QueueUrl: p.client.QueueUrl,
			Entries:  entries,
		})

		if dErr != nil {
			log.Errorf("failed while trying to delete message: %s", dErr)
		}
	}
}
