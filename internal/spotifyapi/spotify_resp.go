package spotifyapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// GetTrackInfo retrieves information about a track from the Spotify API.
func (c *Client) GetTrackInfo(trackURL string) (OriginalTrack, error) {
	req, err := http.NewRequest("GET", trackURL, nil)
	if err != nil {
		return OriginalTrack{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)

	resp, err := c.Client.Do(req)
	if err != nil {
		return OriginalTrack{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return OriginalTrack{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var trackInfo OriginalTrack
	err = json.NewDecoder(resp.Body).Decode(&trackInfo)
	if err != nil {
		return OriginalTrack{}, err
	}

	return trackInfo, nil
}

// GetPlaylist retrieves information about a playlist from the Spotify API.
func (c *Client) GetPlaylist(playlistURL string) (Playlist, error) {
//  https://open.spotify.com/playlist/6uo85AkZ0mkn3rB5P7U6qy?si=a91710ed23ec467b
	req, err := http.NewRequest("GET", playlistURL, nil)
	if err != nil {
		return Playlist{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)

	resp, err := c.Client.Do(req)
	if err != nil {
		return Playlist{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Playlist{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var playlist Playlist
	err = json.NewDecoder(resp.Body).Decode(&playlist)
	if err != nil {
		return Playlist{}, err
	}

	return playlist, nil
}

// ExtractPlaylistID extracts the playlist ID from a Spotify playlist URL.
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
