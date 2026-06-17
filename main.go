package main

import (
	"errors"
	"net"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = "23234"
)

const (
	topbarHeight     = 2
	sidebarFullWidth = 20
	cellWidth        = 3
	cellHeight       = 1
	gutterWidth      = 12
)

func main() {
	srv, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),

		wish.WithHostKeyPath(".ssh/id_ed25519"),

		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),

			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatal("SSHGarden: Could not start server", "error", err)
	}

	log.Info("Starting SSH server", "host", host, "port", port)
	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not start server", "error", err)
	}
}

type coordinate struct {
	x int
	y int
}

type model struct {
	width         int
	height        int
	currentScreen screen
	gardenGrid    [][]rune
	mousePosition coordinate
	gridStartX    int
	gridStartY    int
	sidebarOpen   bool
	selectedPlot  coordinate
	currentTime   time.Time
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
	return tick()
}

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

		coloredTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#2c7450")).Render(title)
		coloredOptions := lipgloss.NewStyle().Foreground(lipgloss.Color("#4484f2")).Render(options)
		fullMenu := lipgloss.JoinVertical(lipgloss.Center, coloredTitle, "\n-------------------------------------\n", coloredOptions)

		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, fullMenu)

	case gardenScreen:
		// Topbar
		clock := m.currentTime.Format("15:04:05")
		leftGutter := lipgloss.NewStyle().Width(gutterWidth).Align(lipgloss.Left).Render(clock)
		rightGutter := lipgloss.NewStyle().Width(gutterWidth).Align(lipgloss.Right).Render("")
		center := lipgloss.NewStyle().Width(m.width - 2*gutterWidth).Align(lipgloss.Center).Render("SSH Garden")
		topBar := lipgloss.JoinHorizontal(lipgloss.Center, leftGutter, center, rightGutter)

		// Sidebar
		var sidebarWidth int

		if m.sidebarOpen {
			sidebarWidth = sidebarFullWidth
		}

		sideBar := "X"

		sidebarHeight := m.height - topbarHeight
		gardenWidth := m.width - sidebarWidth
		gardenHeight := m.height - topbarHeight

		var gridBuilder strings.Builder

		for y := 0; y < len(m.gardenGrid); y++ {
			for x := 0; x < len(m.gardenGrid[y]); x++ {
				gridBuilder.WriteString("[" + string(m.gardenGrid[y][x]) + "]")
			}
			gridBuilder.WriteString("\n")
		}

		grid := gridBuilder.String()

		styledGarden := lipgloss.NewStyle().Width(gardenWidth).Height(gardenHeight).Align(lipgloss.Center, lipgloss.Center).Render(grid)
		styledSidebar := lipgloss.NewStyle().Width(sidebarWidth-1).Height(sidebarHeight).Border(lipgloss.NormalBorder(), false, false, false, true).Align(lipgloss.Right).Render(sideBar)
		styledTopbar := lipgloss.NewStyle().Width(m.width).Border(lipgloss.NormalBorder(), false, false, true, false).Align(lipgloss.Center).Render(topBar)

		centerArea := lipgloss.JoinHorizontal(lipgloss.Top, styledGarden, styledSidebar)
		content = lipgloss.JoinVertical(lipgloss.Left, styledTopbar, centerArea)

	}

	return content
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	grid := [][]rune{}
	for range 5 {
		row := []rune{}
		for range 5 {
			row = append(row, ' ')
		}
		grid = append(grid, row)
	}
	return model{gardenGrid: grid, sidebarOpen: false, currentTime: time.Now()}, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseAllMotion()}
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
