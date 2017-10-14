package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fulhax/mpdbot/ircbot"
	irccmd "github.com/fulhax/mpdbot/ircbot/cmd"
	"github.com/fulhax/mpdbot/mpd"
	"github.com/gorilla/mux"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type mpdbotConfig struct {
	Debug      bool
	Mpd        string
	HttpPort   string
	IrcEnabled bool
	IrcNick    string
	IrcServer  string
}

var (
	queueHandler *mpd.QueueHandler
	mpdClient    *mpd.MpdClient
	config       mpdbotConfig
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

	addr := fmt.Sprintf(":%s", config.HttpPort)

	log.Printf("Listening on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func initConfig() {

	flag.Bool("debug", false, "Enable debug mode")
	flag.String("mpd", "127.0.0.1:6600", "mpd host")
	flag.String("httpPort", "8888", "Http port")
	flag.Bool("ircEnabled", true, "Enable irc bot")
	flag.String("ircNick", "mpdbot", "Irc nick")
	flag.String("ircServer", "127.0.0.1:6697", "irc server")
	flag.Parse()
	viper.BindPFlag("ircNick", flag.Lookup("ircNick"))

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/mpdbot/")
	viper.AddConfigPath("$HOME/.config/mpdbot")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func main() {
	initConfig()

	mpdClient, err := mpd.NewMpdClient(config.Mpd)
	if err != nil {
		fmt.Println(err)
		return
	}

	queueHandler = &mpd.QueueHandler{MpdClient: mpdClient}
	queueHandler.Init()

	if config.IrcEnabled {
		irc := ircbot.New(config.IrcNick, config.IrcServer, true)
		irc.AddCommand(&irccmd.Usage{})
		irc.AddCommand(&IrcMpdNp{mpdClient})
		irc.AddCommand(&IrcAddSong{mpdClient})
	}

	serveApi()
}
