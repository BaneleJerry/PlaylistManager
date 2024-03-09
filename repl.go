package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Playlist Manager > \n")
		for _,cmd := range getCommands(){
			fmt.Printf("-%d: %s\n",cmd.number, cmd.name)
		}
		scanner.Scan()
		text := scanner.Text()
		cleaned := cleanInput(text)
		availableCommands := getCommands()

		if len(cleaned) == 0 {
			continue
		}

		num, err := strconv.Atoi(cleaned[0])
		if err != nil {
			fmt.Println("Invalid command number")
			continue
		}

		var cmd cliCommand
		found := false
		for _, v := range availableCommands {
			if v.number == num {
				cmd = v
				found = true
				break
			}
		}

		if !found {
			fmt.Println("Unknown command number")
			continue
		}


		err = cmd.callback()
		if err != nil {
			fmt.Println(err)
		}
	}
}

type cliCommand struct {
	number      int
	name        string
	description string
	callback    func() error
}

func getCommands() []cliCommand {
	
	return []cliCommand{
		{1, "help", "Displays a help message", commandHelp},
		{2, "exit", "Exit the Pokedex", commandExit},
		{3, "Generate Playlist", "Generates playlists using provided link", generatePlaylist},

	}
}

func cleanInput(str string) []string {
	lowerd := strings.ToLower(str)
	words := strings.Fields(lowerd)
	return words
}

