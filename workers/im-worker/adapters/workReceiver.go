package adapters

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"io"

	"github.com/daniilcdev/insta-magick-bot/workers/im-worker/ports"
	"github.com/daniilcdev/insta-magick-bot/workers/im-worker/types"
)

type WorkReceiver struct {
	W ports.Worker
}

func (wr *WorkReceiver) StartReceiving() {
	http.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		log.Println("new request")

		b, err := io.ReadAll(io.Reader(r.Body))
		defer r.Body.Close()

		log.Println(string(b))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Unmarshal
		var work types.Work
		err = json.Unmarshal(b, &work)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		wr.W.OnWorkReceived(work)
		fmt.Fprintf(w, "scheduled: '%s' : '%s'\n", work.File, work.Filter)
	})

	log.Println("starting http server on port 8008")
	http.ListenAndServe(":8008", nil)
}
