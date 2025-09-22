package listener

import (
	"context"
	"encoding/json"
	"log"
	"sqs-example/internal/app/sqs/processor"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Listener struct {
	eventChannel chan ConsumedEvent
	receiver     Receiver
	processor    Processor
}

type Receiver interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

type Processor interface {
	ProcessMessage(ctx context.Context, m processor.SendingMessage) error
}

func NewListener(r Receiver, p Processor) *Listener {
	return &Listener{
		eventChannel: make(chan ConsumedEvent),
		receiver:     r,
		processor:    p,
	}
}

func (l *Listener) Listen(ctx context.Context, url string) {
	go l.listen(ctx, url)
	go l.process(ctx, url)
}

func (l *Listener) listen(ctx context.Context, url string) {
	for {
		out, err := l.receiver.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &url,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10, // long polling
			VisibilityTimeout:   30, // seconds to process before it reappears
		})
		if err != nil {
			log.Printf("ReceiveMessage error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if len(out.Messages) == 0 {
			continue
		}

		for _, msg := range out.Messages {
			var event ConsumedMessage
			err := json.Unmarshal([]byte(*msg.Body), &event)

			if err != nil {
				log.Printf("Unmarshal error: %v", err)
			}
			l.eventChannel <- ConsumedEvent{
				ReceiptHandle: msg.ReceiptHandle,
				body:          event,
			}
		}
	}
}

func (l *Listener) process(ctx context.Context, url string) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-l.eventChannel:
			l.processEvent(ctx, &e, url)
		}
	}
}

func (l *Listener) processEvent(ctx context.Context, e *ConsumedEvent, url string) {
	err := l.processor.ProcessMessage(ctx, &e.body)

	_, err = l.receiver.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &url,
		ReceiptHandle: e.ReceiptHandle,
	})
	if err != nil {
		log.Printf("Error deleting event: %v", err)
	}
}
