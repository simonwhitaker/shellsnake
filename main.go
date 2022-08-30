package main

import (
	"fmt"
	"math/rand"
	"os"
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
	hasCrashed   bool
	highScore    int
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

func initialModel(highScore int) model {
	return model{
		body:           []coord{{x: 0, y: 0}},
		food:           coord{x: 6, y: 0},
		foodGlyphIndex: 0,
		headings:       []direction{right},
		length:         initLength,
		tickDuration:   time.Millisecond * 150,
		hasCrashed:     false,
		highScore:      highScore,
	}
}

func (m model) score() int {
	return m.length - initLength
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
			// Try to persist high score
			config, _ := LoadConfig()
			config.HighScore = m.highScore
			err := config.Save()
			if err != nil {
				fmt.Printf("Error saving config: %v", err)
			}
			return m, tea.Quit

		case " ":
			if m.hasCrashed {
				m = initialModel(m.highScore)
				return m, tickEvery(m.tickDuration)
			}

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
			m.hasCrashed = true
			return m, nil
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
			m.hasCrashed = true
			return m, nil
		}
		m.body = append(tail, newHead)

		if m.score() > m.highScore {
			m.highScore = m.score()
		}

		return m, tickEvery(m.tickDuration)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	footerStyle := lipgloss.NewStyle().Faint(true).Width(cols * 2).Align(lipgloss.Center)
	crossBar := strings.Repeat("─", cols*2)

	scoreStyle := lipgloss.NewStyle().Bold(true)

	var highScoreStyle lipgloss.Style

	if m.score() < m.highScore {
		highScoreStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#999999"))
	} else {
		highScoreStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00ff00"))
	}

	var headGlyph, bodyGlyph, footerMsg string
	if m.hasCrashed {
		headGlyph = "💀"
		bodyGlyph = "🍖"
		footerMsg = "Space to start, q to quit"
	} else {
		headGlyph = "😄"
		bodyGlyph = "🐛"
		footerMsg = "↑ ↓ ← →, q to quit"
	}

	scoreStr := scoreStyle.Render(fmt.Sprintf(" Score: %d", m.score()))
	highScoreStr := highScoreStyle.Render(fmt.Sprintf("High Score: %d ", m.highScore))
	scoreSpacerWidth := cols*2 - lipgloss.Width(scoreStr) - lipgloss.Width(highScoreStr)

	s := "╭" + crossBar + "╮\n"
	s += "│" + scoreStyle.Render(scoreStr) + strings.Repeat(" ", scoreSpacerWidth) + highScoreStyle.Render(highScoreStr) + "│\n"
	s += "├" + crossBar + "┤\n"

	for r := 0; r < rows; r++ {
		s += "│"
		for c := 0; c < cols; c++ {
			pos := coord{x: c, y: r}
			if pos == m.food {
				s += foodGlyphs[m.foodGlyphIndex]
			} else if pos == m.body[len(m.body)-1] {
				s += headGlyph
			} else if contains(m.body, pos) {
				s += bodyGlyph
			} else {
				s += "  "
			}
		}
		s += "│\n"
	}

	s += "├" + crossBar + "┤\n"
	s += "│" + footerStyle.Render(footerMsg) + "│\n"
	s += "╰" + crossBar + "╯\n"

	// Send the UI for rendering
	return s
}

func main() {
	rand.Seed(time.Now().UnixNano())
	config, _ := LoadConfig()
	p := tea.NewProgram(initialModel(config.HighScore))
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
