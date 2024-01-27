package adapters

import "github.com/daniilcdev/insta-magick-bot/workers/im-worker/types"

type WorkReceiver struct {
	W types.Worker
}

func (receiver *WorkReceiver) ProcessNewFilesInDir(path string) {
	w := types.Work{
		File: path,
	}
	receiver.W.OnWorkReceived(w)
}
