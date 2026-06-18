package main

import (
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
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.title.width, m.title.height = msg.Width, msg.Height
		m.garden.width, m.garden.height = msg.Width, msg.Height
		return m, nil
	case switchToGardenMsg:
		m.currentScreen = gardenScreen
		return m, nil
	case tickMsg:
		var cmd tea.Cmd
		m.garden, cmd = m.garden.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	switch m.currentScreen {
	case titleScreen:
		m.title, cmd = m.title.Update(msg)
	case gardenScreen:
		m.garden, cmd = m.garden.Update(msg)
	}
	return m, cmd
}
