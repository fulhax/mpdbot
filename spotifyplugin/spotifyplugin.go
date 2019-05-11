package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifySong struct {
	client *spotify.Client
}

func (s SpotifySong) IsSong(song string) bool {
	log.Printf("spotify %s", song)
	if len(song) > 31 && song[0:31] == "https://open.spotify.com/track/" {
		return true
	}
	return false
}

func (s *SpotifySong) GetURI(song string) (string, string, error) {
	splits := strings.Split(song, "/")
	if len(splits) != 5 {
		log.Printf("asdasdfas")
		return "", "", errors.New("Failed to parse URI")
	}

	if s.client == nil {
		log.Printf("Spotify client is nil, it will soon crash!")
	}
	log.Printf("split 4 is %s", splits[4])

	track, err := s.client.GetTrack(spotify.ID(splits[4]))
	if err != nil {
		log.Printf("error in spotify %s", err.Error())
		return "", "", err
	}

	log.Printf("Spotify uri - %s", track.URI)
	return string(track.URI), track.Name, nil
}

/**
You need to set the SPOTIFY_ID and SPOTIFY_SECRET environmental variables.
For more info see https://github.com/zmb3/spotify
*/
func (s *SpotifySong) Auth() {
	config := &clientcredentials.Config{
		TokenURL: spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		return
	}

	client := spotify.Authenticator{}.NewClient(token)
	s.client = &client
}

var MpdMediaSource SpotifySong
