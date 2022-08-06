package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const rows = 16
const cols = 16

type coord struct {
	x int
	y int
}

type direction int
type tickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

const (
	up direction = iota
	right
	down
	left
)

type model struct {
	body    []coord
	heading direction
	length  int
}

func initialModel() model {
	return model{
		body:    []coord{{x: 0, y: 0}},
		heading: right,
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
			if m.heading == right || m.heading == left {
				m.heading = up
			}

		case "down":
			if m.heading == right || m.heading == left {
				m.heading = down
			}

		case "left":
			if m.heading == up || m.heading == down {
				m.heading = left
			}

		case "right":
			if m.heading == up || m.heading == down {
				m.heading = right
			}

		}

	case tickMsg:
		current_head := m.body[len(m.body)-1]
		new_head := current_head
		switch m.heading {
		case up:
			new_head.y--
		case down:
			new_head.y++
		case left:
			new_head.x--
		case right:
			new_head.x++
		}
		if new_head.x < 0 || new_head.x >= cols || new_head.y < 0 || new_head.y >= rows {
			return m, tea.Quit
		} else {
			start_index := len(m.body) - m.length + 1
			if start_index < 0 {
				start_index = 0
			}
			tail := m.body[start_index:]
			if contains(tail, new_head) {
				// The snake has collided with itself.
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
	s := "╭" + strings.Repeat("─", cols) + "╮\n"

	for r := 0; r < rows; r++ {
		s += "│"
		for c := 0; c < cols; c++ {
			pos := coord{x: c, y: r}
			if contains(m.body, pos) {
				s += "o"
			} else {
				s += " "
			}
		}
		s += "│\n"
	}

	s += "╰" + strings.Repeat("─", cols) + "╯\n"

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
