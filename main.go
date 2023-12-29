package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/daniilcdev/insta-magick-bot/client/telegram"
	folderscanner "github.com/daniilcdev/insta-magick-bot/workers/folderScanner"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	folderscanner.ProcessFunc = func(entry fs.DirEntry) {
		err := os.Rename("./res/raw/"+entry.Name(), "./res/processed/"+entry.Name())
		if err != nil {
			fmt.Printf("failed to move file %s: %v\n", entry.Name(), err)
		}

		fmt.Println("new file processed")
	}

	go folderscanner.KeepScanning(ctx, "./res/raw", 2*time.Second)

	botClient, err := telegram.NewClassroomTrackerBot("6346977744:AAHePgVewxrkGwZH5KaoVmExzY5wYqddrig")
	if err != nil {
		log.Default().Println(err)
	}

	go botClient.Start()

	fmt.Println("keeping system alive for 10 minutes")
	<-time.After(10 * time.Minute)
}
