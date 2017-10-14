package main

import (
	"log"
	"net/http"
)

type NowPlayingResponse struct {
	State string
	Song  string
}

// TODO: create mpd object with connection and stuff.
func addToPlayerlistHandler(w http.ResponseWriter, r *http.Request) {

}

func status(w http.ResponseWriter, r *http.Request) {
	status, err := mpdClient.GetStatus()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}
	serveJSON(w, r, status)
}

func searchAndAdd(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	user := r.FormValue("user")
	file, err := queueHandler.AddToQueue(user, search)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest)
	}

	serveJSON(w, r, file)
}

func getPlaylist(w http.ResponseWriter, r *http.Request) {
	items, err := mpdClient.GetPlaylist()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	serveJSON(w, r, items)
}

func getNowPlayingHandler(w http.ResponseWriter, r *http.Request) {

	status, err := mpdClient.GetStatus()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	song, err := mpdClient.CurrentSong()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	resp := NowPlayingResponse{
		status.State,
		song,
	}

	serveJSON(w, r, resp)
}

func playNextSongHandler(w http.ResponseWriter, r *http.Request) {
	err := mpdClient.Next()
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	song, err := mpdClient.CurrentSong()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	serveJSON(w, r, song)
}
