package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const Rows = 16
const Cols = 16

type coord struct {
	x int
	y int
}

type Direction int
type TickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

const (
	Up Direction = iota
	Right
	Down
	Left
)

type model struct {
	body    []coord
	heading Direction
	length  int
}

func initialModel() model {
	return model{
		body:    []coord{{x: 0, y: 0}},
		heading: Right,
		length:  10,
	}
}

func (m model) Init() tea.Cmd {
	return tickEvery()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.heading == Right || m.heading == Left {
				m.heading = Up
			}

		case "down":
			if m.heading == Right || m.heading == Left {
				m.heading = Down
			}

		case "left":
			if m.heading == Up || m.heading == Down {
				m.heading = Left
			}

		case "right":
			if m.heading == Up || m.heading == Down {
				m.heading = Right
			}

		}

	case TickMsg:
		current_head := m.body[len(m.body)-1]
		new_head := current_head
		switch m.heading {
		case Up:
			new_head.y--
		case Down:
			new_head.y++
		case Left:
			new_head.x--
		case Right:
			new_head.x++
		}
		if new_head.x < 0 || new_head.x >= Cols || new_head.y < 0 || new_head.y >= Rows {
			return m, tea.Quit
		} else {
			start_index := len(m.body) - m.length + 1
			if start_index < 0 {
				start_index = 0
			}
			tail := m.body[start_index:]
			if contains(tail, new_head) {
				return m, tea.Quit
			}
			m.body = append(tail, new_head)
			return m, tickEvery()
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func contains(s []coord, e coord) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (m model) View() string {
	s := "╭" + strings.Repeat("─", Cols) + "╮\n"

	for r := 0; r < Rows; r++ {
		s += "│"
		for c := 0; c < Cols; c++ {
			pos := coord{x: c, y: r}
			if contains(m.body, pos) {
				s += "o"
			} else {
				s += " "
			}
		}
		s += "│\n"
	}

	s += "╰" + strings.Repeat("─", Cols) + "╯\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
