package telegram

import messaging "github.com/daniilcdev/insta-magick-bot/messaging/pkg"

type WorkScheduler interface {
	Schedule(work messaging.Work) error
}
