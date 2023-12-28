package main

import (
	"log"

	"github.com/daniilcdev/insta-magick-bot/client/telegram"
)

func main() {
	botClient, err := telegram.NewClassroomTrackerBot("6346977744:AAHePgVewxrkGwZH5KaoVmExzY5wYqddrig")
	if err != nil {
		log.Default().Println(err)
	}

	botClient.Start()
}
