package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"

	log "github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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

	// Get TrackNames
	trackNames := []string{os.Args[1], os.Args[2], os.Args[3]}

	// Search all at once so spotify doesn't have to wait
	var waitGroup sync.WaitGroup
	var trackIDs []spotify.ID
	c := make(chan spotify.ID)

	for _, tracks := range trackNames {
		waitGroup.Add(1)
		go searchIDs(client, tracks, c, &waitGroup)
	}

	for range trackNames {
		trackIDs = append(trackIDs, <-c)
	}

	waitGroup.Wait()

	//Build recommend Request
	seeds := spotify.Seeds{
		Artists: []spotify.ID{},
		Tracks:  trackIDs,
		Genres:  []string{},
	}

	//Get Recs from Spotify
	result, err := client.GetRecommendations(seeds, nil, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"seeds":  seeds,
			"result": result,
			"error":  err,
		}).Fatal("Couldn't Get Recs")
	}
	for _, recommendations := range result.Tracks {
		fmt.Println("  ", recommendations.Name, " ", recommendations.Artists)
	}

}

func searchIDs(spot spotify.Client, trackName string, c chan spotify.ID, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	log.WithFields(log.Fields{
		"track_name": trackName,
	}).Debug("Now Processing Search")
	log.Debugf("now processing: %v\n", trackName)
	result, err := spot.Search(trackName, spotify.SearchType(spotify.SearchTypeTrack))
	if err != nil {
		log.WithFields(log.Fields{
			"track_name": trackName,
			"result":     result,
			"error":      err,
		}).Fatal("Failed search")
	}

	log.WithFields(log.Fields{
		"track_name":    trackName,
		"result":        result,
		"total_results": result.Tracks.Total,
	}).Debug("Search Result Retrieved")

	c <- result.Tracks.Tracks[0].ID

}
