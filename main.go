package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	err := tea.NewProgram(&Model{}, tea.WithAltScreen()).Start()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

}

type Model struct {
	count int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return &m, tea.Quit
		case "up":
			{
				m.count = m.count + 1
				return &m, nil
			}
		case "down":
			{
				m.count = m.count - 1
				return &m, nil
			}

		}
	}
	return &m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("count: %d \n\n ↑ increse ↓ decrese", m.count)
}
