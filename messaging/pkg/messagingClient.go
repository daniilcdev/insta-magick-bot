package messaging

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/nats-io/nats.go"
)

type MessagingClient interface {
	io.Closer
	Schedule(work Work) error
	OnMessage(topic string, handler func(context.Context, *Work)) error
}

type natsClient struct {
	ns  *nats.Conn
	ctx context.Context
}

func Connect(ctx context.Context) (MessagingClient, error) {
	ns, err := nats.Connect("nats:4222")
	if err != nil {
		return nil, err
	}

	mq := &natsClient{
		ns:  ns,
		ctx: ctx,
	}

	log.Println("mq connected")

	return mq, nil
}

func (mq *natsClient) Schedule(work Work) error {
	data, err := json.Marshal(work)
	if err != nil {
		log.Printf("failed to submit: '%v'\n", work)
		return err
	}

	if err := mq.ns.Publish(WorkCreated, data); err != nil {
		log.Printf("failed to publish: '%v'\n", work)
		return err
	}

	return nil
}

func (mq *natsClient) OnMessage(topic string, handler func(context.Context, *Work)) error {
	_, err := mq.ns.Subscribe(topic,
		func(msg *nats.Msg) {
			if work, err := getWork(msg); err != nil {
				log.Printf("unable to unmarshal message: '%v'\n", err)
			} else {
				handler(mq.ctx, work)
			}
		},
	)

	return err
}

func (mq *natsClient) Close() error {
	return mq.ns.Drain()
}

func getWork(msg *nats.Msg) (*Work, error) {
	var work Work
	err := json.Unmarshal(msg.Data, &work)
	return &work, err
}
