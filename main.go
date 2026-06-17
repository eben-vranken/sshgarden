package main

import (
	"errors"
	"net"

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

type model struct {
	width         int
	height        int
	currentScreen screen
	gardenGrid    [][]rune
}

type screen int

const (
	titleScreen screen = iota
	gardenScreen
)

func (m model) Init() tea.Cmd {
	return nil
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
		grid := ""

		for y := 0; y < len(m.gardenGrid); y++ {
			for x := 0; x < len(m.gardenGrid[y]); x++ {
				grid += "[" + string(m.gardenGrid[y][x]) + "]"
			}
			grid += "\n"
		}

		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, grid)
	}

	return content
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	grid := [][]rune{}
	for y := 0; y < 5; y++ {
		row := []rune{}
		for x := 0; x < 5; x++ {
			row = append(row, '.')
		}
		grid = append(grid, row)
	}
	return model{gardenGrid: grid}, []tea.ProgramOption{tea.WithAltScreen()}
}
