package main

import "log"

func main() {
	if err := createApp().start(); err != nil {
		log.Fatalln(err)
	}
}
