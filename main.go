package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type coord struct {
	x int
	y int
}

type direction int
type tickMsg time.Time

type model struct {
	body           []coord
	food           coord
	foodGlyphIndex int
	heading        direction
	nextHeading    direction
	length         int
	tickDuration   time.Duration
}

const (
	rows       = 16
	cols       = 16
	initLength = 3
)

const (
	up direction = iota
	right
	down
	left
)

var foodGlyphs = [...]string{"🍌", "🍎", "🍊", "🍐", "🍰"}

func initialModel() model {
	return model{
		body:           []coord{{x: 0, y: 0}},
		food:           coord{x: 6, y: 0},
		foodGlyphIndex: 0,
		heading:        right,
		nextHeading:    right,
		length:         initLength,
		tickDuration:   time.Millisecond * 150,
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

		// Update heading. You can only turn 90 degrees at a time.
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

		if newHead == m.food {
			m.length++
			m.food = getRandomCoord(append(m.body, newHead))
			m.foodGlyphIndex = (m.foodGlyphIndex + 1) % len(foodGlyphs)
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

func (m model) View() string {
	headerStyle := lipgloss.NewStyle().Bold(true).Width(cols * 2).Align(lipgloss.Center)
	footerStyle := lipgloss.NewStyle().Faint(true).Width(cols * 2).Align(lipgloss.Center)

	s := "╭" + strings.Repeat("──", cols) + "╮\n"
	s += "│" + headerStyle.Render("Score: "+strconv.Itoa(m.length-initLength)) + "│\n"
	s += "├" + strings.Repeat("──", cols) + "┤\n"

	for r := 0; r < rows; r++ {
		s += "│"
		for c := 0; c < cols; c++ {
			pos := coord{x: c, y: r}
			if pos == m.food {
				s += foodGlyphs[m.foodGlyphIndex]
			} else if pos == m.body[len(m.body)-1] {
				s += "😄"
			} else if contains(m.body, pos) {
				s += "🐛"
			} else {
				s += "  "
			}
		}
		s += "│\n"
	}

	s += "├" + strings.Repeat("──", cols) + "┤\n"
	s += "│" + footerStyle.Render("↑ ↓ ← →, q to quit") + "│\n"
	s += "╰" + strings.Repeat("──", cols) + "╯\n"

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
