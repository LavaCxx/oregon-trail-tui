package main

import (
	"fmt"
	"os"

	"github.com/LavaCxx/oregon-trail-tui/game"
	"github.com/LavaCxx/oregon-trail-tui/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Set up terminal for UTF-8
	os.Setenv("LANG", "en_US.UTF-8")

	m := game.NewModel()
	p := tea.NewProgram(ui.NewModel(&m), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
