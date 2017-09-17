package main

import (
	"log"
	"net/http"

	"github.com/fhs/gompd/mpd"
)

func getDbQueue() {

}

func getSong(name string) string {
	mpdcon, err := mpd.Dial("tcp", *mpdAddr)
	defer mpdcon.Close()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	// Find why to get all songs from mpd

}

func addToPlaylist(name string) (err error) {
	// items, err := mpdcon.PlaylistInfo(-1, -1)
	// if err != nil {
	// 	log.Fatal(err)
	// 	errorHandler(w, r, http.StatusBadRequest)
	// 	return
	// }
}
