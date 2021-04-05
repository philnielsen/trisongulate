package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	localspot "github.com/philnielsen/trisongulate/spotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

var Tracks []string
var rootCmd = &cobra.Command{
	Use:   "trisongulate",
	Short: "get a list of reccomendations from 3 songs",
	Long:  `Get a list of reccomendations from three songs`,
	Run: func(cmd *cobra.Command, args []string) {
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

		// Get TrackIDs
		trackIDs := localspot.AggregateTrackIDs(client, Tracks)

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
	},
}

func init() {
	rootCmd.Flags().StringSliceVarP(&Tracks, "tracks", "t", []string{}, "Tracks to trisongulate")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
