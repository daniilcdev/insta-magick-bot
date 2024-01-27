package adapters

import (
	"log"
	"os"

	"github.com/daniilcdev/insta-magick-bot/workers/im-worker/ports"
)

type loggingWorkReporter struct {
	logger *log.Logger
}

func NewLoggingReporter() ports.WorkReporter {

	return &loggingWorkReporter{
		logger: log.New(os.Stdout, "im-worker", log.LstdFlags|log.LUTC),
	}
}

func (r *loggingWorkReporter) Done(work any) {
	r.logger.Printf("work '%v' done\n", work)
}

func (r *loggingWorkReporter) Failed(work any) {
	r.logger.Printf("work '%v' failed\n", work)
}
