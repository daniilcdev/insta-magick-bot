package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/daniilcdev/insta-magick-bot/client/telegram"
)

type TriggerHandler interface {
	ProcessNewFilesInDir(dir string, entries []string)
}

type PipelineTrigger struct {
	Handler TriggerHandler
	Storage telegram.Storage
}

func (s *PipelineTrigger) KeepScanning(ctx context.Context, path string, period time.Duration) {
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("pipeline trigger stopped by cancel")
			return
		case <-ticker.C:
			go s.processFiles(path, nil)

			ticker.Reset(period)
		}
	}
}

func (s *PipelineTrigger) processFiles(path string, files []string) {
	if s.Handler == nil {
		fmt.Println("no func of process")
		return
	}

	s.Handler.ProcessNewFilesInDir(path, files)
}
