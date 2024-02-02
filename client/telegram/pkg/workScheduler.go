package telegram

import (
	types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"
)

type WorkScheduler interface {
	Schedule(work types.Work) error
}
