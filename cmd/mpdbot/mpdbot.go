package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fulhax/mpdbot/handler"
	"github.com/fulhax/mpdbot/irccmd"
	"github.com/fulhax/mpdbot/mpd"
	"github.com/fulhax/mpdbot/mpd/statistics"
	"github.com/fulhax/mpdbot/mpd/statistics/sqlite"
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
	StatsDB     string
}

var (
	queueHandler *mpd.QueueHandler
	mpdClient    *mpd.MpdClient
	stats        statistics.Storage
	config       mpdbotConfig
	BuildDate    string
	Version      string
)

func serveHTTP() {
	h := handler.New(mpdClient, queueHandler, stats)
	http.Handle("/", h)

	addr := fmt.Sprintf(":%s", config.HTTPport)
	log.Printf("Listening on %s\n", addr)
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
	flag.String("statsdb", "./stats.db", "statistics database (sqlite)")
	flag.Parse()

	viper.BindPFlag("debug", flag.Lookup("debug"))
	viper.BindPFlag("mpd", flag.Lookup("mpd"))
	viper.BindPFlag("mpdPassword", flag.Lookup("mpdPassword"))
	viper.BindPFlag("httpPort", flag.Lookup("httpPort"))
	viper.BindPFlag("ircServer", flag.Lookup("ircServer"))
	viper.BindPFlag("ircTls", flag.Lookup("ircTls"))
	viper.BindPFlag("ircNick", flag.Lookup("ircNick"))
	viper.BindPFlag("statsDB", flag.Lookup("statsdb"))

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
	if Version != "" && BuildDate != "" {
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage of %s (%s %s):\n", os.Args[0], Version, BuildDate)
			flag.PrintDefaults()
		}
	}

	initConfig()

	var err error
	stats, err = sqlite.New(config.StatsDB)
	if err != nil {
		log.Fatalf("Unable setup stats db (%v)", err.Error())
	}

	mpdClient, err = mpd.NewMpdClient(config.Mpd, config.MpdPassword)
	if err != nil {
		log.Fatalf("Unable to connect to mpd (%v)", err.Error())
		return
	}

	queueHandler = mpd.NewQueueHandler(mpdClient, stats)
	go queueHandler.HandlePlaylist(config.Mpd, config.MpdPassword)

	if config.IrcEnabled {
		irc := ircbot.New(config.IrcNick, config.IrcServer, config.IrcTLS)
		irc.AddCommand(&ircbotcmd.Usage{})
		irc.AddCommand(irccmd.NewNp(mpdClient))
		irc.AddCommand(irccmd.NewAddSong(mpdClient, queueHandler))
		irc.AddCommand(irccmd.NewMpdUpdate(mpdClient))
		irc.AddCommand(irccmd.NewUserTop(stats))
		irc.AddCommand(irccmd.NewAutodj(queueHandler))
	}

	serveHTTP()
}
