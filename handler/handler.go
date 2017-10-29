package handler

import "net/http"
import "log"
import "github.com/gorilla/mux"
import "encoding/json"
import "github.com/fulhax/mpdbot/mpd"

type handler struct {
	mpdClient    *mpd.MpdClient
	queueHandler *mpd.QueueHandler
}

type NowPlayingResponse struct {
	State string
	Song  string
}

func New(m *mpd.MpdClient, q *mpd.QueueHandler) *mux.Router {
	h := handler{m, q}
	r := mux.NewRouter()

	r.HandleFunc("/playlist", jsonResponseHandler(h.getPlaylist)).Methods("GET")
	r.HandleFunc("/current", jsonResponseHandler(h.getNowPlayingHandler)).Methods("GET")
	r.HandleFunc("/next", jsonResponseHandler(h.playNextSongHandler)).Methods("POST")
	r.HandleFunc("/add", jsonResponseHandler(h.searchAndAdd)).Methods("GET")
	r.HandleFunc("/status", jsonResponseHandler(h.status)).Methods("GET")

	return r
}

func jsonResponseHandler(h func(http.ResponseWriter, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}

		w.WriteHeader(status)
		w.Header().Set("content-Type", "application/json")
		err = json.NewEncoder(w).Encode(data)

		if err != nil {
			log.Printf("Could not encode json: %v", err)
		}
	}
}

// TODO: create mpd object with connection and stuff.
func addToPlayerlistHandler(w http.ResponseWriter, r *http.Request) {
}

func (h handler) status(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	status, err := h.mpdClient.GetStatus()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return status, http.StatusOK, nil
}

func (h handler) searchAndAdd(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	search := r.FormValue("search")
	user := r.FormValue("user")
	file, err := h.queueHandler.AddToQueue(user, search)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return file, http.StatusOK, nil
}

func (h handler) getPlaylist(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	items, err := h.mpdClient.GetPlaylist()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return items, http.StatusOK, nil
}

func (h handler) getNowPlayingHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {

	status, err := h.mpdClient.GetStatus()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	song, err := h.mpdClient.CurrentSong()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	resp := NowPlayingResponse{
		status.State,
		song,
	}

	return resp, http.StatusOK, nil
}

func (h handler) playNextSongHandler(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	err := h.mpdClient.Next()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	song, err := h.mpdClient.CurrentSong()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return song, http.StatusOK, nil
}
