package telegram

import (
	types "github.com/daniilcdev/insta-magick-bot/image-service-worker/pkg"
)

type WorkScheduler interface {
	Schedule(work types.Work) error
}
