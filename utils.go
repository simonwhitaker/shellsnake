package main

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func isHorizontal(d direction) bool {
	return d == left || d == right
}

func getRandomCoord(exclude []coord) coord {
	for {
		x := rand.Intn(cols)
		y := rand.Intn(rows)
		if !contains(exclude, coord{x: x, y: y}) {
			return coord{x: x, y: y}
		}
	}
}

func contains(s []coord, e coord) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
