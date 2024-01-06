package folderscanner

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"time"
)

type FileHandler interface {
	ProcessNewFile(dir string, entries []fs.DirEntry)
}

type FolderScanner struct {
	FoundFileHandler FileHandler
}

func (s *FolderScanner) KeepScanning(ctx context.Context, path string, period time.Duration) {
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("folder scanner stopped by cancel")
			return
		case <-ticker.C:
			nFiles, err := os.ReadDir(path)

			switch {
			case err != nil:
				fmt.Println(err)
			case len(nFiles) > 0:
				const nWorkers = 4
				cap := min(len(nFiles), nWorkers)
				go s.processFiles(path, nFiles[:cap])
			}

			ticker.Reset(period)
		}
	}
}

func (s *FolderScanner) processFiles(path string, files []fs.DirEntry) {
	if s.FoundFileHandler == nil {
		fmt.Println("no func of process")
		return
	}
	s.FoundFileHandler.ProcessNewFile(path, files)
}
