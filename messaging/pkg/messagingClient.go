package messaging

import (
	"encoding/json"
	"io"
	"log"

	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"
	"github.com/nats-io/nats.go"
)

type MessagingClient interface {
	io.Closer
	Schedule(work types.Work) error
	Notify(topic string, ch chan *types.Work)
}

type natsClient struct {
	ns   *nats.Conn
	subs map[string][]chan *types.Work
}

func InitMessageQueue() MessagingClient {
	ns, err := nats.Connect("nats:4222")
	if err != nil {
		log.Fatalln(err)
	}

	mq := &natsClient{
		ns:   ns,
		subs: make(map[string][]chan *types.Work),
	}

	log.Println("mq initialized")

	return mq
}

func (mq *natsClient) Schedule(work types.Work) error {
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

func (mq *natsClient) Notify(topic string, receiver chan *types.Work) {
	if channels, ok := mq.subs[topic]; ok {
		mq.subs[topic] = append(channels, receiver)
	} else {
		channels = make([]chan *types.Work, 0, 4)
		mq.subs[topic] = append(channels, receiver)
		mq.ns.Subscribe(topic, mq.handleMessage)
	}
}

func (mq *natsClient) Close() error {
	return mq.ns.Drain()
}

func (mq *natsClient) handleMessage(msg *nats.Msg) {
	defer msg.Ack()

	c, ok := mq.subs[msg.Subject]
	if !ok {
		log.Printf("no handlers for topic '%s'", msg.Subject)
		return
	}

	work, err := getWork(msg)
	if err != nil {
		log.Printf("unable to unmarshal work: '%v'\n", err)
		return
	}

	for _, ch := range c {
		ch <- work
	}
}

func getWork(msg *nats.Msg) (*types.Work, error) {
	var work types.Work
	err := json.Unmarshal(msg.Data, &work)
	return &work, err
}
