package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fulhax/mpdbot/handler"
	"github.com/fulhax/mpdbot/irccmd"
	"github.com/fulhax/mpdbot/mpd"
	"github.com/rendom/ircbot"
	ircbotcmd "github.com/rendom/ircbot/cmd"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type mpdbotConfig struct {
	Debug       bool
	Mpd         string
	MpdPassword string
	HTTPport    string
	IrcEnabled  bool
	IrcNick     string
	IrcTLS      bool
	IrcServer   string
}

var (
	queueHandler *mpd.QueueHandler
	mpdClient    *mpd.MpdClient
	config       mpdbotConfig
)

func serveHTTP() {
	h := handler.New(mpdClient, queueHandler)
	http.Handle("/", h)

	addr := fmt.Sprintf(":%s", config.HTTPport)
	log.Printf("Listening on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func initConfig() {

	flag.Bool("debug", false, "Enable debug mode")
	flag.String("mpd", "127.0.0.1:6600", "mpd host")
	flag.String("mpdPassword", "", "mpd password")
	flag.String("httpPort", "8888", "Http port")
	flag.Bool("ircEnabled", true, "Enable irc bot")
	flag.String("ircNick", "mpdbot", "Irc nick")
	flag.String("ircServer", "127.0.0.1:6697", "irc server")
	flag.Bool("ircTls", true, "irc tls")
	flag.Parse()

	viper.BindPFlag("debug", flag.Lookup("debug"))
	viper.BindPFlag("mpd", flag.Lookup("mpd"))
	viper.BindPFlag("mpdPassword", flag.Lookup("mpdPassword"))
	viper.BindPFlag("httpPort", flag.Lookup("httpPort"))
	viper.BindPFlag("ircServer", flag.Lookup("ircServer"))
	viper.BindPFlag("ircTls", flag.Lookup("ircTls"))
	viper.BindPFlag("ircNick", flag.Lookup("ircNick"))

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/mpdbot/")
	viper.AddConfigPath("$HOME/.config/mpdbot")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initConfig()

	var err error
	mpdClient, err = mpd.NewMpdClient(config.Mpd, config.MpdPassword)
	if err != nil {
		log.Fatal(err)
		return
	}

	queueHandler = &mpd.QueueHandler{MpdClient: mpdClient}
	queueHandler.Init()

	if config.IrcEnabled {
		irc := ircbot.New(config.IrcNick, config.IrcServer, config.IrcTLS)
		irc.AddCommand(&ircbotcmd.Usage{})
		irc.AddCommand(irccmd.NewNp(mpdClient))
		irc.AddCommand(irccmd.NewAddSong(mpdClient, queueHandler))
		irc.AddCommand(irccmd.NewMpdUpdate(mpdClient))
	}

	serveHTTP()
}
