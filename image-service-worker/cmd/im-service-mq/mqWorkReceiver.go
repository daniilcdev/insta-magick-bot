package mq

import (
	"encoding/json"
	"log"

	types "github.com/daniilcdev/insta-magick-bot/image-service-worker/pkg"
	messaging "github.com/daniilcdev/insta-magick-bot/messaging/pkg"

	"github.com/nats-io/nats.go"
)

type MQWorkReceiver struct {
	W  Worker
	nc *nats.Conn
}

func (wr *MQWorkReceiver) StartReceiving() {
	var nc *nats.Conn
	var err error
	if nc, err = nats.Connect("nats:4222"); err != nil {
		log.Fatalln(err)
		return
	}

	wr.nc = nc
	if _, err = nc.QueueSubscribe(messaging.WorkCreated, "workers", wr.onWorkCreated); err != nil {
		log.Printf("failed to queue-subscribe: '%v'\n", err)
		return
	}

	log.Println("worker subscribed")
}

func (wr *MQWorkReceiver) Close() {
	wr.nc.Close()
	wr.W = nil
}

func (wr *MQWorkReceiver) onWorkCreated(msg *nats.Msg) {
	var work types.Work
	if err := json.Unmarshal(msg.Data, &work); err != nil {
		return
	}

	if err := wr.W.Do(work); err != nil {
		log.Printf("work failed: '%v'\n", err)
		wr.failed(work)
		msg.Nak()
		return
	}

	wr.done(work)
	msg.Ack()
}

func (wr *MQWorkReceiver) done(work types.Work) {
	data, err := json.Marshal(work)
	if err != nil {
		log.Printf("failed to serialize failed work: '%v'\n", work)
		return
	}

	wr.nc.Publish(messaging.WorkDone, data)
}

func (wr *MQWorkReceiver) failed(work types.Work) {
	data, err := json.Marshal(work)
	if err != nil {
		log.Printf("failed to serialize failed work: '%v'\n", work)
		return
	}

	wr.nc.Publish(messaging.WorkFailed, data)
}
