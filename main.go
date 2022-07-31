package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	err := tea.NewProgram(initializeModel(), tea.WithAltScreen()).Start()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

}

type Model struct {
	cursor   int
	options  []string
	selected map[int]struct{}
	spinner  spinner.Model
}

func (m Model) Init() tea.Cmd {
	return spinner.Tick
}

func initializeModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Moon
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Model{
		cursor:   0,
		options:  []string{"Counter", "Wheather"},
		selected: make(map[int]struct{}),
		spinner:  s,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return &m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	}
	return &m, nil
}

func (m Model) View() string {
	// return fmt.Sprintf("count: %d \n\n ↑ increse ↓ decrese", m.count)
	s := "What operation you would like to perform \n\n"

	for i, option := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		result := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
			result = m.spinner.View()
		}

		s += fmt.Sprintf("%s [%s] %s %s \n", cursor, checked, option, result)
	}

	s += "\nPress q to quit.\n"

	return s
}

// func getWheather() (w string) {
// 	var client = &http.Client{
// 		Timeout: 10 * time.Second,
// 	}
// 	res, err := client.Get("https://charm.sh/")
// 	if err != nil {
// 		panic("API failed")
// 	}
// 	return res.Status
// }
