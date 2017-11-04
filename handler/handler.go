package handler

import "net/http"
import "log"
import "github.com/gorilla/mux"
import "encoding/json"
import "github.com/fulhax/mpdbot/mpd"
import "github.com/fulhax/mpdbot/mpd/statistics"

type handler struct {
	mpdClient    *mpd.MpdClient
	queueHandler *mpd.QueueHandler
	stats        statistics.Storage
}

type NowPlayingResponse struct {
	State string
	Song  string
}

func New(m *mpd.MpdClient, q *mpd.QueueHandler, s statistics.Storage) *mux.Router {
	h := handler{m, q, s}
	r := mux.NewRouter()

	r.HandleFunc("/current", jsonResponseHandler(h.getNowPlaying)).Methods("GET")
	r.HandleFunc("/next", jsonResponseHandler(h.playNextSong)).Methods("POST")
	r.HandleFunc("/add", jsonResponseHandler(h.searchAndAdd)).Methods("GET")
	r.HandleFunc("/search", jsonResponseHandler(h.search)).Methods("GET")
	r.HandleFunc("/status", jsonResponseHandler(h.status)).Methods("GET")
	r.HandleFunc("/top", jsonResponseHandler(h.toplist)).Methods("GET")

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

func (h handler) search(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	search := r.FormValue("search")
	result, err := h.mpdClient.SearchInLibrary(search)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return result, http.StatusOK, nil
}

func (h handler) searchAndAdd(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	search := r.FormValue("search")
	user := r.FormValue("user")

	result, err := h.mpdClient.SearchInLibrary(search)
	if err != nil || len(result) == 0 {
		return nil, http.StatusBadRequest, err
	}

	file, err := h.queueHandler.AddToQueue(user, result[0].File, result[0].File)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return file, http.StatusOK, nil
}

func (h handler) getNowPlaying(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
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

func (h handler) playNextSong(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
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

func (h handler) toplist(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	user := r.FormValue("search")

	var items []statistics.SongStats
	var err error

	if user != "" {
		items, err = h.stats.GetUserTop(user, 25)
	} else {
		items, err = h.stats.GetTop(25)
	}

	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return items, http.StatusOK, nil
}
