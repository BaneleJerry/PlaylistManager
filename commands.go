package main

import (
	"bufio"
	"fmt"
	"os"

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
	playlistmanager.GeneratePlaylist(playListURL)
	return nil
}
