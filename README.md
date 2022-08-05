# trisongulate

Take three songs and recommend a playlist of songs.

PHIL ADDED THIS

# Running UI

1. `go install ./...`
1. fill in a `local.env` with Spotify creds
1. `trisongulate server`
1. navigate to `http://localhost:8080/` in your browser
1. Search for tunes!

### To Run Headless

1. `go install ./...`
1. fill in local.env with Spotify creds
1. `trisongulate -t "ramble on" -t "baba o'riley" -t "we will rock you"`
1. Enjoy those tunes that are spit out!

#### Awesome Libraries Used

- https://github.com/zmb3/spotify
