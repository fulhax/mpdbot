package main

import (
	"log"
	"strings"
)

type YouTube struct {
}

func (s YouTube) isSong(song string, q *QueueHandler) bool {
	log.Printf("youtube, %s", song)
	splits := strings.Split(song, ".")
	if len(splits) == 3 && strings.ToLower(splits[1]) == "youtube" {
		return true
	}
	return false
}

//TODO: Make sure that the url is http and not https
func (s YouTube) getURI(song string, q *QueueHandler) string {
	return "yt:" + song
}
