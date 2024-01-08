package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/daniilcdev/insta-magick-bot/client/telegram"
)

type TriggerHandler interface {
	ProcessNewFilesInDir(dir string, entries []string)
}

type PipelineTrigger struct {
	Handler TriggerHandler
	storage telegram.Storage
}

func (s *PipelineTrigger) KeepScanning(ctx context.Context, path string, period time.Duration) {
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("folder scanner stopped by cancel")
			return
		case <-ticker.C:
			const batchSize = 10

			pending := s.storage.GetPendingRequests(batchSize)

			if len(pending) > 0 {
				log.Println("begin IM processing pipeline")
				go s.processFiles(path, pending)
			}

			ticker.Reset(period)
		}
	}
}

func (s *PipelineTrigger) processFiles(path string, files []string) {
	if s.FoundFilesHandler == nil {
		fmt.Println("no func of process")
		return
	}

	s.FoundFilesHandler.ProcessNewFilesInDir(path, files)
}
