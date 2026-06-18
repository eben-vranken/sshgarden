package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type gardenModel struct {
	width         int
	height        int
	gardenGrid    [][]rune
	mousePosition coordinate
	sidebarOpen   bool
	selectedPlot  coordinate
	currentTime   time.Time
}

func (g gardenModel) computeLayout() layout {
	var sidebarWidth int

	if g.sidebarOpen {
		sidebarWidth = sidebarFullWidth
	}

	gardenWidth := g.width - sidebarWidth
	gardenHeight := g.height - topbarHeight
	sidebarHeight := g.height - topbarHeight
	gridStartX := (gardenWidth - len(g.gardenGrid[0])*cellWidth) / 2
	gridStartY := topbarHeight + (gardenHeight-len(g.gardenGrid)*cellHeight)/2

	return layout{
		sidebarWidth:  sidebarWidth,
		gardenWidth:   gardenWidth,
		gardenHeight:  gardenHeight,
		sidebarHeight: sidebarHeight,
		gridStartX:    gridStartX,
		gridStartY:    gridStartY,
	}
}

func (g gardenModel) cellAt(px, py int) (coordinate, bool) {
	l := g.computeLayout()
	col := (px - l.gridStartX) / cellWidth
	row := (py - l.gridStartY) / cellHeight

	inbounds := (px >= l.gridStartX && py >= l.gridStartY) && (col >= 0 && col < len(g.gardenGrid[0]) && row >= 0 && row < len(g.gardenGrid))

	return coordinate{x: col, y: row}, inbounds
}

func (g gardenModel) View() string {
	l := g.computeLayout()

	// Topbar
	clock := g.currentTime.Format("15:04:05")
	leftGutter := lipgloss.NewStyle().Width(gutterWidth).Align(lipgloss.Left).Foreground(colorBeige).Render(clock)
	rightGutter := lipgloss.NewStyle().Width(gutterWidth).Align(lipgloss.Right).Render("")
	center := lipgloss.NewStyle().Width(g.width - 2*gutterWidth).Align(lipgloss.Center).Foreground(colorLightGreen).Bold(true).Render("SSH Garden")
	topBar := lipgloss.JoinHorizontal(lipgloss.Center, leftGutter, center, rightGutter)

	sideBar := lipgloss.NewStyle().Foreground(colorRedOrange).Bold(true).Render("X")

	var gridBuilder strings.Builder

	for y := 0; y < len(g.gardenGrid); y++ {
		if y > 0 {
			gridBuilder.WriteString("\n")
		}
		for x := 0; x < len(g.gardenGrid[y]); x++ {
			gridBuilder.WriteString("[" + string(g.gardenGrid[y][x]) + "]")
		}
	}

	grid := gridBuilder.String()

	styledGarden := lipgloss.NewStyle().Width(l.gardenWidth).Height(l.gardenHeight).Align(lipgloss.Center, lipgloss.Center).Foreground(colorGrayGreen).Render(grid)
	styledTopbar := lipgloss.NewStyle().Width(g.width).Align(lipgloss.Center).Render(topBar)

	var separator string
	var centerArea string
	if g.sidebarOpen {
		separator = strings.Repeat("═", l.gardenWidth) + "╦" + strings.Repeat("═", g.width-l.gardenWidth-1)
		styledSidebar := lipgloss.NewStyle().Width(l.sidebarWidth-1).Height(l.sidebarHeight).Border(lipgloss.DoubleBorder(), false, false, false, true).BorderForeground(colorGray).Align(lipgloss.Right).Render(sideBar)
		centerArea = lipgloss.JoinHorizontal(lipgloss.Top, styledGarden, styledSidebar)
	} else {
		separator = strings.Repeat("═", g.width)
		centerArea = styledGarden
	}

	separator = lipgloss.NewStyle().Foreground(colorGray).Render(separator)

	return lipgloss.JoinVertical(lipgloss.Left, styledTopbar, separator, centerArea)
}

func (g gardenModel) Update(msg tea.Msg) (gardenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:

		switch msg.Action {
		case tea.MouseActionMotion:
			g.mousePosition, _ = g.cellAt(msg.X, msg.Y)
		case tea.MouseActionPress:
			if msg.Button == tea.MouseButtonLeft {
				if g.sidebarOpen && msg.X == g.width-1 && msg.Y == topbarHeight {
					g.sidebarOpen = false
				} else if cell, ok := g.cellAt(msg.X, msg.Y); ok {
					g.sidebarOpen = true
					g.selectedPlot = cell
				}
			}
		}
	case tickMsg:
		g.currentTime = time.Time(msg)
		return g, tick()
	}

	return g, nil
}

func (g gardenModel) Init() tea.Cmd {
	return tick()
}
