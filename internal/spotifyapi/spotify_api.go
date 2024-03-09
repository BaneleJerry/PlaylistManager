package spotifyapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const baseURL string = "https://api.spotify.com/v1"
type Client struct {
	http.Client
	token tokenResp
}



func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func NewClient() *Client {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Please set SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET in the .env file")
	}

	client := &Client{
		Client: http.Client{
			Timeout: time.Minute / 2,
		},
	}

	token, err := client.GetToken(clientID, clientSecret)
	if err != nil {
		log.Fatalf("Error getting token: %v", err)
	}

	client.token = token
	return client
}

func (c *Client) GetToken(clientID, clientSecret string) (tokenResp, error) {
	body := url.Values{}
	body.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token",
		strings.NewReader(body.Encode()))

	if err != nil {
		return tokenResp{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := clientID + ":" + clientSecret
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))

	resp, err := c.Client.Do(req)
	if err != nil {
		return tokenResp{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return tokenResp{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return tokenResp{}, err
	}

	tokenResp := tokenResp{}

	err = json.Unmarshal(dat, &tokenResp)
	if err != nil {
		return tokenResp, err
	}

	return tokenResp, nil
}
