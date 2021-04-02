package handlers

import (
	"html/template"
	"net/http"

	localspot "github.com/philnielsen/trisongulate/spotify"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
)

func Home(w http.ResponseWriter, r *http.Request, spot localspot.Trisongulate) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	ts, err := template.ParseFS(content, "templates/home.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// We then use the Execute() method on the template set to write the template
	// content as the response body. The last parameter to Execute() represents any
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	// client := localspot.GetClient(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func Search(w http.ResponseWriter, r *http.Request, spot localspot.Trisongulate) {
	if r.URL.Path != "/search" {
		http.NotFound(w, r)
		return
	}

	r.ParseForm()
	var trackNames []string
	for _, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			trackNames = append(trackNames, value)
		}
	}

	log.WithFields(log.Fields{
		"trackNames": trackNames,
	}).Info("Using Seeds")

	// Get TrackIDs
	trackIDs := localspot.AggregateTrackIDs(spot.SpotifyClient, trackNames)

	//Build recommend Request
	seeds := spotify.Seeds{
		Artists: []spotify.ID{},
		Tracks:  trackIDs,
		Genres:  []string{},
	}

	log.WithFields(log.Fields{
		"seeds": seeds,
	}).Info("Using Seeds")

	//Get Recs from Spotify
	result, err := spot.SpotifyClient.GetRecommendations(seeds, nil, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"seeds":  seeds,
			"result": result,
			"error":  err,
		}).Fatal("Couldn't Get Recs")
	}

	var fullTracks []*spotify.FullTrack
	for _, rec := range result.Tracks {
		fullTrack, err := spot.SpotifyClient.GetTrack(rec.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"name":   rec.Name,
				"result": fullTrack,
				"error":  err,
			}).Info("Couldn't Get Track")
		}

		fullTracks = append(fullTracks, fullTrack)

	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	ts, err := template.ParseFS(content, "templates/results.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// We then use the Execute() method on the template set to write the template
	// content as the response body. The last parameter to Execute() represents any
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	// client := localspot.GetClient(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))

	err = ts.Execute(w, fullTracks)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
