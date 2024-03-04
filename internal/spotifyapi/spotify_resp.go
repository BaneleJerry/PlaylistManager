package spotifyapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func (c *Client) GetToken() (tokenResp, error) {
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

func (c *Client) GetTrackInfo(accessToken string) (map[string]interface{}, error) {
	// Create request
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/tracks/4cOdK2wGLETKBW3PvgPWqT", nil)
	if err != nil {
		return nil, err
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response body
	var trackInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&trackInfo)
	if err != nil {
		return nil, err
	}

	// Return track info
	return trackInfo, nil
}

func (c *Client) GetPlaylist(accessToken string) (Playlist, error) {
	// endpoint := baseURl + "playlist/"
	// playlistID:=
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/6uo85AkZ0mkn3rB5P7U6qy", nil)
	if err != nil {
		return Playlist{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Send request
	resp, err := c.Client.Do(req)
	if err != nil {
		return Playlist{}, err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return Playlist{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response body
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Playlist{}, err
	}

	playlist := Playlist{}
	err = json.Unmarshal(dat, &playlist)
	if err != nil {
		return playlist, err
	}

	// Return track info
	return playlist, nil
}

func ExtractPlaylistID(url string) string {
	// Define a regular expression pattern to match the playlist ID
	pattern := regexp.MustCompile(`playlist/(\w+)`)

	// Find the first match in the URL
	matches := pattern.FindStringSubmatch(url)

	// Extract the playlist ID from the matched string
	if len(matches) > 1 {
		return matches[1]
	}

	// Return an empty string if no match found
	return ""
}
