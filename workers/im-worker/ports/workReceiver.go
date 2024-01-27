package ports

import "github.com/daniilcdev/insta-magick-bot/workers/im-worker/types"

type Worker interface {
	OnWorkReceived(work types.Work)
}
