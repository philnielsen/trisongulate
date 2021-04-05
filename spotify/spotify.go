package spotify

import "github.com/zmb3/spotify"

type Trisongulate struct {
	SpotifyClient spotify.Client
}

type Reccomendations struct {
	Tracks []spotify.FullTrack
}
