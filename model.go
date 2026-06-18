package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type coordinate struct {
	x int
	y int
}

type model struct {
	width         int
	height        int
	currentScreen screen
	title         titleModel
	garden        gardenModel
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type screen int

const (
	titleScreen screen = iota
	gardenScreen
)

func (m model) Init() tea.Cmd {
	return tea.Batch(m.title.Init(), m.garden.Init())
}
