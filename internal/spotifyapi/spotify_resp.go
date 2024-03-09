package spotifyapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
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
func (c *Client) getPlaylist(playlistURL string) (Playlist, error) {
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

func (c *Client) getNextPage(PlaylistURL string) (getPlaylistItemsRESP, error) {
	req, err := http.NewRequest("GET", PlaylistURL, nil)
	if err != nil {
		return getPlaylistItemsRESP{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)

	resp, err := c.Client.Do(req)
	if err != nil {
		return getPlaylistItemsRESP{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getPlaylistItemsRESP{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var playlist getPlaylistItemsRESP
	err = json.NewDecoder(resp.Body).Decode(&playlist)
	if err != nil {
		return getPlaylistItemsRESP{}, err
	}

	return playlist, nil
}

func (c *Client) GetSongs(playlistURL string) ([]Track, string, error) {
	var allTracks []Track

	// Initial request to fetch the first set of songs
	playlist, err := c.getPlaylist(playlistURL)
	if err != nil {
		return nil, "", err
	}
	allTracks = append(allTracks, extractTracksFromPlaylist(&playlist)...)
	playlistName := playlist.Name
	// Fetch next sets of songs until there are no more
	nextPage := *playlist.Tracks.Next
	for nextPage != "" {
		// Make request to fetch next set of songs
		playlist, err := c.getNextPage(nextPage)
		if err != nil {
			return nil, "", err
		}
		allTracks = append(allTracks, extractTracksFromPlaylistRESP(&playlist)...)
		nextPage = playlist.Next
	}

	return allTracks, playlistName, nil
}

func extractTracksFromPlaylist(playlist *Playlist) []Track {
	var tracks []Track
	for _, item := range playlist.Tracks.Items {
		var artists []string
		for _, artist := range item.Track.Artists {
			artists = append(artists, artist.Name)
		}

		tracks = append(tracks, Track{
			Title:  item.Track.Name,
			Artist: strings.Join(artists, "/"),
			Path:   "",
		})
	}
	return tracks
}

func extractTracksFromPlaylistRESP(playlist *getPlaylistItemsRESP) []Track {
	var tracks []Track
	for _, item := range playlist.Items {
		var artists []string
		for _, artist := range item.Track.Artists {
			artists = append(artists, artist.Name)
		}

		tracks = append(tracks, Track{
			Title:  item.Track.Name,
			Artist: strings.Join(artists, "/"),
			Path:   "",
		})
	}
	return tracks
}

// ExtractPlaylistID extracts the playlist ID from a Spotify playlist URL.
func extractPlaylistID(url string) string {
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
