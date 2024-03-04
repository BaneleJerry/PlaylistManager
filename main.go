package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	playlistmanager "github.com/BaneleJerry/PlaylistManager/internal/playlistManager"
	"github.com/BaneleJerry/PlaylistManager/internal/spotifyapi"
	"github.com/dhowden/tag"
)

func main() {

	client := spotifyapi.NewClient()

	tokenResp, err := client.GetToken()
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}


	if tokenResp.AccessToken == "" {
		log.Fatal("Access token is empty")
	}

	spotifyPlaylist, err := client.GetPlaylist(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("Error getting track information: %v", err)
	}


	spotifyTracks := extractTracksFromPlaylist(spotifyPlaylist)


	directory := "testMusicFolder"
	localTracks, err := readLocalMusicFiles(directory)
	if err != nil {
		log.Fatalf("Error reading local music files: %v", err)
	}

	var playlistTracks []spotifyapi.Track


	uniqueTracks := make(map[string]bool)

	for _, spspotifyTrack := range spotifyTracks {
		trackID := fmt.Sprintf("%s-%s", spspotifyTrack.Title, spspotifyTrack.Artist)

		if uniqueTracks[trackID] {
			continue
		}
		uniqueTracks[trackID] = true

		for _, localTrack := range localTracks {
			if spspotifyTrack.Title == localTrack.Title && spspotifyTrack.Artist == localTrack.Artist {
				// Append the track to the playlistTracks slice
				playlistTracks = append(playlistTracks, localTrack)
				break 
			}
		}
	}

	fmt.Println(spotifyPlaylist.Name)
	err = playlistmanager.CreateM3U(playlistTracks, spotifyPlaylist.Name)
	if err != nil {
		log.Fatalf("Error creating M3U playlist: %v", err)
	}
}


// extractTracksFromPlaylist extracts track information from Spotify playlist
func extractTracksFromPlaylist(playlist spotifyapi.Playlist) []spotifyapi.Track {
	var tracks []spotifyapi.Track
	for _, item := range playlist.Tracks.Items {
		var artists []string
		for _, artist := range item.Track.Artists {
			artists = append(artists, artist.Name)
		}
		tracks = append(tracks, spotifyapi.Track{
			Title:  item.Track.Name,
			Artist: strings.Join(artists, ", "),
			Path:   "",
		})
	}
	return tracks
}

// readLocalMusicFiles reads local music files from the specified directory
func readLocalMusicFiles(directory string) ([]spotifyapi.Track, error) {
	var localTracks []spotifyapi.Track

	// Read directory contents
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	// Iterate over files in directory
	for _, file := range files {
		// Construct the full file path
		filePath := filepath.Join(directory, file.Name())

		// Check if the item is a file (not a directory)
		if !file.IsDir() && isMusicFile(filePath) {
			// Get the absolute path of the file
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				log.Printf("Error getting absolute path for file %s: %v", filePath, err)
				continue
			}

			// Open the file
			f, err := os.Open(absPath)
			if err != nil {
				log.Printf("Error opening file %s: %v", absPath, err)
				continue
			}
			defer f.Close() // Close the file when done

			// Extract metadata from music file
			metadata, err := tag.ReadFrom(f)
			if err != nil {
				log.Printf("Error reading metadata from %s: %v", absPath, err)
				continue
			}

			// Create Track object with full path and append to localTracks slice
			localTracks = append(localTracks, spotifyapi.Track{
				Title:  metadata.Title(),
				Artist: metadata.Artist(),
				Path:   absPath,
			})
		}
	}
	return localTracks, nil
}

// isMusicFile checks if a file has a music file extension
func isMusicFile(path string) bool {
	musicFileExtensions := []string{".mp3", ".wav", ".flac"}
	ext := strings.ToLower(filepath.Ext(path))
	for _, musicExt := range musicFileExtensions {
		if ext == musicExt {
			return true
		}
	}
	return false
}
