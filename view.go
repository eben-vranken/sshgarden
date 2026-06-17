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

type layout struct {
	sidebarWidth  int
	gardenWidth   int
	gardenHeight  int
	sidebarHeight int
	gridStartX    int
	gridStartY    int
}

func (m model) computeLayout() layout {
	var sidebarWidth int

	if m.sidebarOpen {
		sidebarWidth = sidebarFullWidth
	}

	gardenWidth := m.width - sidebarWidth
	gardenHeight := m.height - topbarHeight
	sidebarHeight := m.height - topbarHeight
	gridStartX := (gardenWidth - len(m.gardenGrid[0])*cellWidth) / 2
	gridStartY := topbarHeight + (gardenHeight-len(m.gardenGrid)*cellHeight)/2

	return layout{
		sidebarWidth:  sidebarWidth,
		gardenWidth:   gardenWidth,
		gardenHeight:  gardenHeight,
		sidebarHeight: sidebarHeight,
		gridStartX:    gridStartX,
		gridStartY:    gridStartY,
	}
}

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
		l := m.computeLayout()

		// Topbar
		clock := m.currentTime.Format("15:04:05")
		leftGutter := lipgloss.NewStyle().Width(gutterWidth).Align(lipgloss.Left).Foreground(colorBeige).Render(clock)
		rightGutter := lipgloss.NewStyle().Width(gutterWidth).Align(lipgloss.Right).Render("")
		center := lipgloss.NewStyle().Width(m.width - 2*gutterWidth).Align(lipgloss.Center).Foreground(colorLightGreen).Bold(true).Render("SSH Garden")
		topBar := lipgloss.JoinHorizontal(lipgloss.Center, leftGutter, center, rightGutter)

		sideBar := lipgloss.NewStyle().Foreground(colorRedOrange).Bold(true).Render("X")

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

		styledGarden := lipgloss.NewStyle().Width(l.gardenWidth).Height(l.gardenHeight).Align(lipgloss.Center, lipgloss.Center).Foreground(colorGrayGreen).Render(grid)
		styledTopbar := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(topBar)

		var separator string
		var centerArea string
		if m.sidebarOpen {
			separator = strings.Repeat("═", l.gardenWidth) + "╦" + strings.Repeat("═", m.width-l.gardenWidth-1)
			styledSidebar := lipgloss.NewStyle().Width(l.sidebarWidth-1).Height(l.sidebarHeight).Border(lipgloss.DoubleBorder(), false, false, false, true).BorderForeground(colorGray).Align(lipgloss.Right).Render(sideBar)
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
