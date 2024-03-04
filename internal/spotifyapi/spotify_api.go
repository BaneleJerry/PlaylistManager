package spotifyapi

import (
    "net/http"
    "os"
    "time"
)

const baseURl = "https://api.spotify.com/v1/"

var clientID = os.Getenv("SPOTIFY_CLIENT_ID")
var clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

func init() {
    if clientID == "" || clientSecret == "" {
        panic("Please set SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET")
    }
}

type Client struct {
    Client http.Client
}

func NewClient() Client {
    return Client{
        Client: http.Client{
            Timeout: time.Minute / 2,
        },
    }
}
