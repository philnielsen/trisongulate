package spotify

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
)

func AggregateTrackIDs(spot spotify.Client, trackNames []string) []spotify.ID {
	// Search all at once so spotify doesn't have to wait
	var waitGroup sync.WaitGroup
	var trackIDs []spotify.ID
	c := make(chan spotify.ID)

	for _, tracks := range trackNames {
		waitGroup.Add(1)
		go searchIDs(spot, tracks, c, &waitGroup)
	}

	for range trackNames {
		trackIDs = append(trackIDs, <-c)
	}

	waitGroup.Wait()
	return trackIDs
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
