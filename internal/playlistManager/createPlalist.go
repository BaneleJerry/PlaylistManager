package playlistmanager

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/BaneleJerry/PlaylistManager/internal/spotifyapi"
)

var osType = runtime.GOOS
var osMusicDir string

const appFolder string = "PlaylistManager"

func init() {
	if osType == "linux" {
		currentUser, err := user.Current()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		osMusicDir = currentUser.HomeDir + "/Music"
	} else if osType == "windows" {
		userProfile := os.Getenv("USERPROFILE")
		osMusicDir = filepath.Join(userProfile, "Music")
	}
}

func CreateM3U(tracks []spotifyapi.Track, playlistName string) error {
    filePath := filepath.Join(osMusicDir, playlistName+".m3u")
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write M3U header
    if _, err := fmt.Fprintf(file, "#EXTM3U\n"); err != nil {
        return err
    }

    // Write track entries
    for _, track := range tracks {
        if _, err := fmt.Fprintf(file, "#EXTINF:%s - %s\n%s\n", track.Artist, track.Title, track.Path); err != nil {
            return err
        }
    }

    fmt.Printf("Playlist '%s' created successfully!\n", playlistName)
    return nil
}


func GetMusicDir() string {
	return osMusicDir
}
