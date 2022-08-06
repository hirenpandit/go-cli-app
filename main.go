package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	score1   string
	score2   string
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
		options:  []string{"Score"},
		selected: make(map[int]struct{}),
		spinner:  s,
		score1:   "",
		score2:   "",
	}
}

type TickMsg time.Time

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
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
		case "r":
			return &m, tickEvery()
		}
	case TickMsg:
		score := getScore()
		fmt.Println("-->", score)
		if len(score) > 0 {
			m.score1 = fmt.Sprintf("%s R: %f / %f O: %f \n", score[0].inning, score[0].r, score[0].w, score[0].o)
			m.score2 = fmt.Sprintf("%s R: %f / %f O: %f \n", score[1].inning, score[1].r, score[1].w, score[1].o)
		}
		return &m, tickEvery()
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return &m, cmd

	}
	return &m, cmd
}

func tickEvery() tea.Cmd {
	return tea.Every(time.Millisecond*3000, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) View() string {
	return optionView(&m)

}

func optionView(m *Model) string {
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

		s += fmt.Sprintf(
			"%s [%s] %s %s \n %s\n %s\n",
			cursor,
			checked,
			option,
			result,
			m.score1,
			m.score2,
		)
	}

	s += "\nPress r to refresh.\n"
	s += "\nPress q to quit.\n"

	return s
}

type Score struct {
	r      float64
	w      float64
	o      float64
	inning string
}

func getScore() []Score {

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	res, err := client.Get("https://api.cricapi.com/v1/currentMatches?apikey=e2e0863d-5e3d-4669-8d5b-d30c3a5017c0&offset=0")
	if err != nil {
		panic("Failed to get score")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic("Failed to read response")
	}
	m := map[string]interface{}{}
	json.Unmarshal(body, &m)
	if data, ok := m["data"]; ok {
		if list, found := data.([]interface{}); found {
			for i := range list {
				mat, _ := json.Marshal(list[i])
				match := map[string]interface{}{}
				json.Unmarshal(mat, &match)
				if mt, k := match["name"].(string); k {
					if mt == "West Indies vs India, 4th T20I" {
						scores, _ := match["score"].([]interface{})
						score := []Score{}
						for s := range scores {
							inning, _ := json.Marshal(scores[s])
							inn := map[string]interface{}{}
							json.Unmarshal(inning, &inn)
							score = append(score, Score{
								r:      inn["r"].(float64),
								o:      inn["o"].(float64),
								w:      inn["w"].(float64),
								inning: inn["inning"].(string),
							})
						}
						return score
					}
				}
			}
		}
	}
	return nil
}
