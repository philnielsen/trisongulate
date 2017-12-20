package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	//Get Auth Token for Spotify
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)

	//Get Track IDs
	trackIDs := []spotify.ID{spotify.ID(os.Args[1]), spotify.ID(os.Args[2]), spotify.ID(os.Args[3])}

	//Build recommend Request
	seeds := spotify.Seeds{
		Artists: []spotify.ID{},
		Tracks:  trackIDs,
		Genres:  []string{},
	}

	//Get Recs from Spotify
	recs, err := client.GetRecommendations(seeds, nil, nil)
	if err != nil {
		log.Fatalf("couldn't get Recs: %v", err)
	}
	for _, recommendations := range recs.Tracks {
		fmt.Println("  ", recommendations.Name, " ", recommendations.Artists)
	}

}
