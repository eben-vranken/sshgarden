package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		case "s":
			m.currentScreen = gardenScreen
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.MouseMsg:
		if m.currentScreen != gardenScreen {
			break
		}

		switch msg.Action {
		case tea.MouseActionMotion:
			m.mousePosition, _ = m.cellAt(msg.X, msg.Y)
		case tea.MouseActionPress:
			if msg.Button == tea.MouseButtonLeft {
				if m.sidebarOpen && msg.X == m.width-1 && msg.Y == topbarHeight {
					m.sidebarOpen = false
				} else if cell, ok := m.cellAt(msg.X, msg.Y); ok {
					m.sidebarOpen = true
					m.selectedPlot = cell
				}
			}
		}
	case tickMsg:
		m.currentTime = time.Time(msg)
		return m, tick()
	}

	return m, nil
}

func (m model) cellAt(px, py int) (coordinate, bool) {
	l := m.computeLayout()
	col := (px - l.gridStartX) / cellWidth
	row := (py - l.gridStartY) / cellHeight

	inbounds := (px >= l.gridStartX && py >= l.gridStartY) && (col >= 0 && col < len(m.gardenGrid[0]) && row >= 0 && row < len(m.gardenGrid))

	return coordinate{x: col, y: row}, inbounds
}
