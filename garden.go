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
	gardenGrid    [][]plot
	mousePosition coordinate
	selectedPlot  coordinate
	currentTime   time.Time
	activeAction  action
}

type action int

const (
	actionNone action = iota
	actionPlant
	actionWater
	actionHarvest
	actionShop
)

func (g gardenModel) computeLayout() layout {
	gardenWidth := g.width
	gardenHeight := g.height - topbarHeight - bottomBarHeight
	gridStartX := (gardenWidth - len(g.gardenGrid[0])*cellWidth) / 2
	gridStartY := topbarHeight + (gardenHeight-len(g.gardenGrid)*cellHeight)/2

	return layout{
		gardenWidth:  gardenWidth,
		gardenHeight: gardenHeight,
		gridStartX:   gridStartX,
		gridStartY:   gridStartY,
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

	var gridBuilder strings.Builder

	for y := 0; y < len(g.gardenGrid); y++ {
		if y > 0 {
			gridBuilder.WriteString("\n\n")
		}
		for x := 0; x < len(g.gardenGrid[y]); x++ {
			currentPlant := g.gardenGrid[y][x]
			gridBuilder.WriteString(lipgloss.NewStyle().Foreground(colorGrayGreen).Render(" ["))

			if currentPlant.plant != nil {
				gridBuilder.WriteString(lipgloss.NewStyle().Foreground(currentPlant.plant.color).Render(string(currentPlant.getGrowthVisual())))
			} else {
				gridBuilder.WriteString(lipgloss.NewStyle().Render(string(currentPlant.getGrowthVisual())))
			}

			gridBuilder.WriteString(lipgloss.NewStyle().Foreground(colorGrayGreen).Render("] "))
		}
	}

	grid := gridBuilder.String()

	// Bottom bar
	bottomBar := lipgloss.NewStyle().Width(g.width).Height(bottomBarHeight-1).Align(lipgloss.Center).Border(lipgloss.DoubleBorder(), true, false, false, false).Render("")

	styledGarden := lipgloss.NewStyle().Width(l.gardenWidth).Height(l.gardenHeight).Align(lipgloss.Center, lipgloss.Center).Render(grid)
	styledTopbar := lipgloss.NewStyle().Width(g.width).Align(lipgloss.Center).Border(lipgloss.DoubleBorder(), false, false, true, false).Render(topBar)

	return lipgloss.JoinVertical(lipgloss.Left, styledTopbar, styledGarden, bottomBar)
}

func (g gardenModel) Update(msg tea.Msg) (gardenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:

		switch msg.Action {
		case tea.MouseActionMotion:
			g.mousePosition, _ = g.cellAt(msg.X, msg.Y)
		case tea.MouseActionPress:
			if msg.Button == tea.MouseButtonLeft {
				// Click
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
