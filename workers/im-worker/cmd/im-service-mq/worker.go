package mq

import types "github.com/daniilcdev/insta-magick-bot/workers/im-worker/pkg"

type Worker interface {
	Do(work types.Work) error
}
