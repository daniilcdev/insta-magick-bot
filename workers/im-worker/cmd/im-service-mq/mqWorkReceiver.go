package mq

import (
	"encoding/json"
	"log"

	"github.com/daniilcdev/insta-magick-bot/internal"
	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"

	"github.com/nats-io/nats.go"
)

type MQWorkReceiver struct {
	W  Worker
	nc *nats.Conn
}

func (wr *MQWorkReceiver) StartReceiving() {
	var nc *nats.Conn
	var err error
	if nc, err = nats.Connect(nats.DefaultURL); err != nil {
		log.Fatalln(err)
		return
	}

	wr.nc = nc
	_, err = nc.Subscribe(internal.WorkCreated, wr.onWorkCreated)
	if err != nil {
		log.Fatalf("failed to subscribe: topic '%s'\n", err)
	}

	log.Println("worker subscriped")
}

func (wr *MQWorkReceiver) Close() {
	wr.nc.Close()
	wr.W = nil
}

func (wr *MQWorkReceiver) onWorkCreated(msg *nats.Msg) {
	log.Println(string(msg.Data))

	var work types.Work
	if err := json.Unmarshal(msg.Data, &work); err != nil {
		return
	}

	if err := wr.W.Do(work); err != nil {
		wr.failed(work)
		return
	}

	wr.done(work)
}

func (wr *MQWorkReceiver) done(work any) {
	data, err := json.Marshal(work)
	if err != nil {
		log.Printf("failed to serialize failed work: '%v'\n", work)
		return
	}

	wr.nc.Publish(internal.WorkDone, data)
}

func (wr *MQWorkReceiver) failed(work any) {
	data, err := json.Marshal(work)
	if err != nil {
		log.Printf("failed to serialize failed work: '%v'\n", work)
		return
	}

	wr.nc.Publish(internal.WorkFailed, data)
}
