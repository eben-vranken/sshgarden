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

		m.recomputeGrid()

	case tea.MouseMsg:
		if m.currentScreen != gardenScreen {
			break
		}

		switch msg.Action {
		case tea.MouseActionMotion:
			m.mousePosition.x = (msg.X - m.gridStartX) / cellWidth
			m.mousePosition.y = (msg.Y - m.gridStartY) / cellHeight
		case tea.MouseActionPress:
			if msg.Button == tea.MouseButtonLeft {
				if m.sidebarOpen && msg.X == m.width-1 && msg.Y == topbarHeight {
					m.sidebarOpen = false
					m.recomputeGrid()
				} else {
					col := (msg.X - m.gridStartX) / cellWidth
					row := (msg.Y - m.gridStartY) / cellHeight
					if msg.X >= m.gridStartX && msg.Y >= m.gridStartY {
						if col >= 0 && col < len(m.gardenGrid[0]) && row >= 0 && row < len(m.gardenGrid) {
							m.sidebarOpen = true
							m.selectedPlot = coordinate{
								x: (msg.X - m.gridStartX) / cellWidth,
								y: (msg.Y - m.gridStartY) / cellHeight,
							}
							m.recomputeGrid()
						}
					}
				}
			}
		}
	case tickMsg:
		m.currentTime = time.Time(msg)
		return m, tick()
	}

	return m, nil
}

func (m *model) recomputeGrid() {
	var sidebarWidth int

	if m.sidebarOpen {
		sidebarWidth = sidebarFullWidth
	}

	gardenWidth := m.width - sidebarWidth
	gardenHeight := m.height - topbarHeight

	m.gridStartX = (gardenWidth - len(m.gardenGrid[0])*cellWidth) / 2
	m.gridStartY = topbarHeight + (gardenHeight-len(m.gardenGrid)*cellHeight)/2
}
