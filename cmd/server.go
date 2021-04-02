package cmd

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"github.com/philnielsen/trisongulate/handlers"
	localspot "github.com/philnielsen/trisongulate/spotify"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var t localspot.Trisongulate

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run the ui server for trisongulate",
	Long:  `run the ui server for trisongulate`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load("local.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		t.GetClientCredntialsClient(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))
		// albums, _ := t.SpotifyClient.NewReleases()
		// log.Println(albums.Albums)
		// first start an HTTP server
		// http.HandleFunc("/callback", localspot.CompleteAuth)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("hi")
			handlers.Home(w, r, t)
		})

		http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
			log.Println("hello")
			handlers.Search(w, r, t)
		})

		// client := localspot.GetClient(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))
		// log.Println("Got request for:", r.URL.String())
		// user, err := client.CurrentUser()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Fprintf(w, "Hello: %s", user.DisplayName)
		http.ListenAndServe(":8080", nil)

		// use the client to make calls that require authorizatio
	},
}
