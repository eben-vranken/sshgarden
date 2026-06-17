package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	topbarHeight     = 2
	sidebarFullWidth = 20
	cellWidth        = 3
	cellHeight       = 1
	gutterWidth      = 12
)

func (m model) View() string {
	var content string
	switch m.currentScreen {
	case titleScreen:
		title := `
$$$$$$\   $$$$$$\  $$\   $$\          $$$$$$\                            $$\                     
$$  __$$\ $$  __$$\ $$ |  $$ |        $$  __$$\                           $$ |                    
$$ /  \__|$$ /  \__|$$ |  $$ |        $$ /  \__| $$$$$$\   $$$$$$\   $$$$$$$ | $$$$$$\  $$$$$$$\  
\$$$$$$\  \$$$$$$\  $$$$$$$$ |$$$$$$\ $$ |$$$$\  \____$$\ $$  __$$\ $$  __$$ |$$  __$$\ $$  __$$\ 
 \____$$\  \____$$\ $$  __$$ |\______|$$ |\_$$ | $$$$$$$ |$$ |  \__|$$ /  $$ |$$$$$$$$ |$$ |  $$ |
$$\   $$ |$$\   $$ |$$ |  $$ |        $$ |  $$ |$$  __$$ |$$ |      $$ |  $$ |$$   ____|$$ |  $$ |
\$$$$$$  |\$$$$$$  |$$ |  $$ |        \$$$$$$  |\$$$$$$$ |$$ |      \$$$$$$$ |\$$$$$$$\ $$ |  $$ |
 \______/  \______/ \__|  \__|         \______/  \_______|\__|       \_______| \_______|\__|  \__|`

		options := `[q] Quit [s] Start`

		coloredTitle := lipgloss.NewStyle().Foreground(colorDarkGreen).Render(title)
		coloredOptions := lipgloss.NewStyle().Foreground(colorLightBlue).Render(options)
		fullMenu := lipgloss.JoinVertical(lipgloss.Center, coloredTitle, "\n-------------------------------------\n", coloredOptions)

		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, fullMenu)

	case gardenScreen:
		// Topbar
		clock := m.currentTime.Format("15:04:05")
		leftGutter := lipgloss.NewStyle().Width(gutterWidth).Align(lipgloss.Left).Foreground(colorBeige).Render(clock)
		rightGutter := lipgloss.NewStyle().Width(gutterWidth).Align(lipgloss.Right).Render("")
		center := lipgloss.NewStyle().Width(m.width - 2*gutterWidth).Align(lipgloss.Center).Foreground(colorLightGreen).Bold(true).Render("SSH Garden")
		topBar := lipgloss.JoinHorizontal(lipgloss.Center, leftGutter, center, rightGutter)

		// Sidebar
		var sidebarWidth int

		if m.sidebarOpen {
			sidebarWidth = sidebarFullWidth
		}

		sideBar := lipgloss.NewStyle().Foreground(colorRedOrange).Bold(true).Render("X")

		sidebarHeight := m.height - topbarHeight
		gardenWidth := m.width - sidebarWidth
		gardenHeight := m.height - topbarHeight

		var gridBuilder strings.Builder

		for y := 0; y < len(m.gardenGrid); y++ {
			if y > 0 {
				gridBuilder.WriteString("\n")
			}
			for x := 0; x < len(m.gardenGrid[y]); x++ {
				gridBuilder.WriteString("[" + string(m.gardenGrid[y][x]) + "]")
			}
		}

		grid := gridBuilder.String()

		styledGarden := lipgloss.NewStyle().Width(gardenWidth).Height(gardenHeight).Align(lipgloss.Center, lipgloss.Center).Foreground(colorGrayGreen).Render(grid)
		styledTopbar := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(topBar)

		var separator string
		var centerArea string
		if m.sidebarOpen {
			separator = strings.Repeat("═", gardenWidth) + "╦" + strings.Repeat("═", m.width-gardenWidth-1)
			styledSidebar := lipgloss.NewStyle().Width(sidebarWidth-1).Height(sidebarHeight).Border(lipgloss.DoubleBorder(), false, false, false, true).BorderForeground(colorGray).Align(lipgloss.Right).Render(sideBar)
			centerArea = lipgloss.JoinHorizontal(lipgloss.Top, styledGarden, styledSidebar)
		} else {
			separator = strings.Repeat("═", m.width)
			centerArea = styledGarden
		}

		separator = lipgloss.NewStyle().Foreground(colorGray).Render(separator)

		content = lipgloss.JoinVertical(lipgloss.Left, styledTopbar, separator, centerArea)

	}

	return content
}
