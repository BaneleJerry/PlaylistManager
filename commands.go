package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"

	playlistmanager "github.com/BaneleJerry/PlaylistManager/internal/playlistManager"
)

func commandHelp() error {

	availableCommands := getCommands()
	fmt.Println("")
	fmt.Println("Welcome PLaylist Manger \n Usage:")
	for _, cmd := range availableCommands {

		fmt.Printf("-%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println("")
	return nil
}

func commandExit() error {
	fmt.Println("Bye Bye!!!")
	os.Exit(0)

	return nil
}

func generatePlaylist() error {
	fmt.Println("Please provide link to  Spotify playlist \nexample(https://open.spotify.com/playlist/37i9dQZF1DZ06evO30JqGN?si=2d070d29194a461d)")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	playListURL := scanner.Text()
	err := playlistmanager.GeneratePlaylist(playListURL)
	if err != nil {
		return err
	}
	return nil
}

func downloadPlaylist() error {
	fmt.Println("Please provide link to  Spotify playlist \nexample(https://open.spotify.com/playlist/37i9dQZF1DZ06evO30JqGN?si=2d070d29194a461d)")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	playListURL := scanner.Text()

	if playListURL == "" {
		return errors.New("you didnt provide a playlist link")
	}

	cmd := exec.Command("python", "-m", "spotdl", playListURL)

	dir := "/home/banele/Workspace/github.com/BaneleJerry/Playlist Manager/testMusicFolder"
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	if err := playlistmanager.GeneratePlaylist(playListURL); err != nil {
		return err
	}

	return nil
}
