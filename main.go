package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	// dbug
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("failed to open log file:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// run
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas tea has failed me: %v\n", err)
		os.Exit(1)
	}

}
