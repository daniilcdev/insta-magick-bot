package mq

import types "github.com/daniilcdev/insta-magick-bot/image-service-worker/pkg"

type Worker interface {
	Do(work types.Work) error
}
