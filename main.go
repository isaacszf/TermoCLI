package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"isaacszf.termooo.net/termooo"
)

func main() {
	p := tea.NewProgram(termooo.InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
