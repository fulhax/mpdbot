package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	debug        *bool   = flag.Bool("debug", false, "Enable debug mode")
	port         *string = flag.String("port", "8888", "Serve api on port")
	mpdAddr      *string = flag.String("mpd", "127.0.0.1:6600", "Mpd")
	dbFile       *string = flag.String("db", "mpdapi.db", "Path to database file")
	queueHandler *QueueHandler
	mpdClient    *MpdClient
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	log.Printf("errorHandler status: %d", status)
	switch status {
	case 404:
		fmt.Fprint(w, "404")
	}
}

func serveJSON(w http.ResponseWriter, r *http.Request, val interface{}) {
	w.Header().Set("content-Type", "application/json")
	b, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func serveApi() {
	r := mux.NewRouter()
	r.HandleFunc("/playlist", getPlaylist).Methods("GET")
	r.HandleFunc("/current", getNowPlayingHandler).Methods("GET")
	r.HandleFunc("/next", playNextSongHandler).Methods("POST")
	r.HandleFunc("/add", searchAndAdd).Methods("GET")
	r.HandleFunc("/status", status).Methods("GET")
	//r.HandleFunc("/search", searchInLibrary).Methods("GET")
	// r.HandleFunct("/add", addToPlaylistHandler).Methods("POST")
	http.Handle("/", r)

	addr := fmt.Sprintf(":%s", *port)

	log.Printf("Listening on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	flag.Parse()
	mpdClient, err := NewMpdClient("127.0.0.1:6600")
	if err != nil {
		fmt.Println(err)
		return
	}
	queueHandler = &QueueHandler{mpdClient: mpdClient}
	queueHandler.Init()

	serveApi()
}
