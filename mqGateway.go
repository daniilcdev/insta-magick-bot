package main

import (
	"encoding/json"
	"log"

	"github.com/daniilcdev/insta-magick-bot/internal"
	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"
	"github.com/nats-io/nats.go"
)

type MQGateway struct {
	ns *nats.Conn
}

var mq *MQGateway

func InitMessageQueue() {
	ns, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln(err)
	}

	ns.Subscribe(internal.WorkFailed, onFailed)
	ns.Subscribe(internal.WorkDone, onDone)

	mq = &MQGateway{
		ns: ns,
	}

	log.Println("mq initialized")
}

func (mq *MQGateway) Schedule(work types.Work) error {
	data, err := json.Marshal(work)
	if err != nil {
		log.Printf("failed to submit: '%v'\n", work)
		return err
	}

	if err := mq.ns.Publish(internal.WorkCreated, data); err != nil {
		log.Printf("failed to publish: '%v'\n", work)
		return err
	}

	return nil
}

func onFailed(msg *nats.Msg) {
	work, err := getWork(msg)

	if err != nil {
		log.Printf("unable to unmarshal on fail: '%v'\n", err)
		return
	}

	log.Printf("failed work: '%v'\n", work)
}

func onDone(msg *nats.Msg) {
	work, err := getWork(msg)

	if err != nil {
		log.Printf("unable to unmarshal on success: '%v'\n", err)
		return
	}

	log.Printf("successful work: '%v'\n", work)
}

func getWork(msg *nats.Msg) (*types.Work, error) {
	var work types.Work
	err := json.Unmarshal(msg.Data, &work)
	return &work, err
}
