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
	// The coordinates of the body of the snake, ordered from tail to head.
	body []coord

	// The coordinates of the current food.
	food coord

	// The index of the current food glyph.
	foodGlyphIndex int

	// The heading(s) of the snake. As keys are pressed, we append to this list,
	// and "play back" the changes one tick at a time.
	headings []direction

	// The current length of the snake
	length       int
	tickDuration time.Duration
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
		headings:       []direction{right},
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
			m.headings = appendHeadingIfLegal(m.headings, up)

		case "down":
			m.headings = appendHeadingIfLegal(m.headings, down)

		case "left":
			m.headings = appendHeadingIfLegal(m.headings, left)

		case "right":
			m.headings = appendHeadingIfLegal(m.headings, right)

		}

	case tickMsg:
		newHead := m.body[len(m.body)-1]

		if len(m.headings) > 1 {
			m.headings = m.headings[1:]
		}
		heading := m.headings[0]

		switch heading {
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
	crossBar := strings.Repeat("─", cols*2)

	s := "╭" + crossBar + "╮\n"
	s += "│" + headerStyle.Render("Score: "+strconv.Itoa(m.length-initLength)) + "│\n"
	s += "├" + crossBar + "┤\n"

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

	s += "├" + crossBar + "┤\n"
	s += "│" + footerStyle.Render("↑ ↓ ← →, q to quit") + "│\n"
	s += "╰" + crossBar + "╯\n"

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
