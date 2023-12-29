package folderscanner

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"time"
)

var ProcessFunc func(fs.DirEntry)

func KeepScanning(ctx context.Context, path string, period time.Duration) {
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
				go processFiles(nFiles[:cap])
			}

			ticker.Reset(period)
		}
	}
}

func processFiles(files []fs.DirEntry) {
	if ProcessFunc == nil {
		fmt.Println("no func of process")
		return
	}

	for _, de := range files {
		ProcessFunc(de)
	}
}
