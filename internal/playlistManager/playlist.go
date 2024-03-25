package playlistmanager

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BaneleJerry/PlaylistManager/internal/spotifyapi"
	"github.com/dhowden/tag"
)

func GeneratePlaylist(playlistUrl string) error {
	client := spotifyapi.NewClient()
	if client == nil {
		return fmt.Errorf("client not assigned")
	}
	spotifyTracks, playlistname, err := client.GetSongs(playlistUrl)
	if err != nil {
		return fmt.Errorf("error getting track information: %v", err)
	}

	directory := "testMusicFolder"
	localTracks, err := readLocalMusicFiles(directory)
	if err != nil {
		return fmt.Errorf("error reading local music files: %v", err)
	}

	var playlistTracks []spotifyapi.Track

	uniqueTracks := make(map[string]bool)

	for _, spspotifyTrack := range spotifyTracks {
		trackID := spspotifyTrack.Title
		if uniqueTracks[trackID] {
			continue
		}
		uniqueTracks[trackID] = true

		for _, localTrack := range localTracks {
			if spspotifyTrack.Title == localTrack.Title && spspotifyTrack.Artist == localTrack.Artist {
				playlistTracks = append(playlistTracks, localTrack)
				break
			}
		}
	}

	fmt.Println(playlistname)
	err = CreateM3U(playlistTracks, playlistname)
	if err != nil {
		return fmt.Errorf("error creating M3U playlist: %v", err)
	}
	return nil
}

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
