package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const rows = 16
const cols = 16

type coord struct {
	x int
	y int
}

type direction int
type tickMsg time.Time

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

const (
	up direction = iota
	right
	down
	left
)

func isHorizontal(d direction) bool {
	return d == left || d == right
}

type model struct {
	body         []coord
	heading      direction
	nextHeading  direction
	length       int
	tickDuration time.Duration
}

func initialModel() model {
	return model{
		body:         []coord{{x: 0, y: 0}},
		heading:      right,
		nextHeading:  right,
		length:       5,
		tickDuration: time.Millisecond * 150,
	}
}

func (m model) Init() tea.Cmd {
	return tickEvery(m.tickDuration)
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
			m.nextHeading = up

		case "down":
			m.nextHeading = down

		case "left":
			m.nextHeading = left

		case "right":
			m.nextHeading = right

		}

	case tickMsg:
		newHead := m.body[len(m.body)-1]

		// Update heading
		if isHorizontal(m.heading) != isHorizontal(m.nextHeading) {
			m.heading = m.nextHeading
		}

		switch m.heading {
		case up:
			newHead.y--
		case down:
			newHead.y++
		case left:
			newHead.x--
		case right:
			newHead.x++
		}

		if newHead.x < 0 || newHead.x >= cols || newHead.y < 0 || newHead.y >= rows {
			return m, tea.Quit
		}

		tailStartIndex := len(m.body) - m.length + 1
		if tailStartIndex < 0 {
			tailStartIndex = 0
		}
		tail := m.body[tailStartIndex:]
		if contains(tail, newHead) {
			// The snake has collided with itself.
			return m, tea.Quit
		}
		m.body = append(tail, newHead)
		return m, tickEvery(m.tickDuration)
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
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#ccff88"))
	s := "╭─" + strings.Repeat("──", cols) + "╮\n"

	for r := 0; r < rows; r++ {
		s += "│ "
		for c := 0; c < cols; c++ {
			pos := coord{x: c, y: r}
			if contains(m.body, pos) {
				s += style.Render("o ")
			} else {
				s += "  "
			}
		}
		s += "│\n"
	}

	s += "├─" + strings.Repeat("──", cols) + "┤\n"
	instructions := "→ ← ↓ ↑, q to quit"
	width := cols * 2
	fmtString := fmt.Sprintf("│ %%-%ds│\n", width)
	s += fmt.Sprintf(fmtString, instructions)
	// s += "│ " + instructions + strings.Repeat("  ", cols-len(instructions)/2) + "│\n"
	s += "╰─" + strings.Repeat("──", cols) + "╯\n"

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
