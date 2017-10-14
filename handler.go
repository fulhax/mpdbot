package main

import (
	"log"
	"net/http"

	"github.com/fhs/gompd/mpd"
)

type Song struct {
	Artist string
	Song   string
	Album  string
	Date   string
	Genre  string
}
type NowPlayingResponse struct {
	State string
	Song  Song
}

// TODO: create mpd object with connection and stuff.
func addToPlayerlistHandler(w http.ResponseWriter, r *http.Request) {

}

func status(w http.ResponseWriter, r *http.Request) {
	mpdcon, err := mpd.Dial("tcp", config.Mpd)
	defer mpdcon.Close()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}
	status, err := mpdcon.Status()
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
	mpdcon, err := mpd.Dial("tcp", config.Mpd)
	defer mpdcon.Close()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	items, err := mpdcon.PlaylistInfo(-1, -1)
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	serveJSON(w, r, items)
}

func getNowPlayingHandler(w http.ResponseWriter, r *http.Request) {
	mpdcon, err := mpd.Dial("tcp", config.Mpd)
	defer mpdcon.Close()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	status, err := mpdcon.Status()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	song, err := mpdcon.CurrentSong()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	resp := NowPlayingResponse{
		status["state"],
		Song{
			song["Artist"],
			song["Title"],
			song["Album"],
			song["Date"],
			song["Genre"],
		},
	}

	serveJSON(w, r, resp)
}

func playNextSongHandler(w http.ResponseWriter, r *http.Request) {
	mpdcon, err := mpd.Dial("tcp", config.Mpd)
	defer mpdcon.Close()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	err = mpdcon.Next()
	if err != nil {
		// 	log.Fatal(err)
		// 	errorHandler(w, r, http.StatusBadRequest)
		return
	}

	song, err := mpdcon.CurrentSong()
	if err != nil {
		log.Fatal(err)
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	resp := Song{
		song["Artist"],
		song["Title"],
		song["Album"],
		song["Date"],
		song["Genre"],
	}

	serveJSON(w, r, resp)
}
