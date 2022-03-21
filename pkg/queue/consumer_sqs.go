package queue

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type ConsumerSQS struct {
	client          *Client
	messagesChannel chan []*sqs.Message
	receiver        ReceiverSQS
}

func NewConsumerSQS(client *Client) *ConsumerSQS {
	c := make(chan []*sqs.Message)
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	return &ConsumerSQS{
		client:          client,
		messagesChannel: c,
		receiver: ReceiverSQS{client: client, messagesChannel: c,
			shutdown: shutdown,
		},
	}

}

func (c *ConsumerSQS) Consume(handler func(m *sqs.Message) error) {
	log.Printf("Starting to consume messages from: %s", *c.client.QueueUrl)
	c.startReceivers()
	c.startProcessor(handler)
}

func (c *ConsumerSQS) DefaultHandler(m *sqs.Message) error {
	log.Infof("message received: %s", *(m.Body))

	return nil
}

func (c *ConsumerSQS) startReceivers() {
	for i := 0; i < c.client.ConsumerConfig.Receivers; i++ {
		go c.receiver.receiveMessages()
	}
}

func (c *ConsumerSQS) startProcessor(handler func(m *sqs.Message) error) {
	p := ProcessorSQS{
		client:  c.client,
		handler: handler,
	}

	for messages := range c.messagesChannel {
		go p.processMessages(messages)
	}
}
