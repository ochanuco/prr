package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ochanuco/prr/internal/tui"
)

func main() {
	repo := flag.String("repo", "", "GitHub repository in owner/repo format")
	state := flag.String("state", "open", "PR state: open, closed, or all")
	flag.Parse()

	m := tui.New(*repo, *state)
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
