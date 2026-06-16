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
	width  int
	height int
}

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
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	title := `
$$$$$$\   $$$$$$\  $$\   $$\          $$$$$$\                            $$\                     
$$  __$$\ $$  __$$\ $$ |  $$ |        $$  __$$\                           $$ |                    
$$ /  \__|$$ /  \__|$$ |  $$ |        $$ /  \__| $$$$$$\   $$$$$$\   $$$$$$$ | $$$$$$\  $$$$$$$\  
\$$$$$$\  \$$$$$$\  $$$$$$$$ |$$$$$$\ $$ |$$$$\  \____$$\ $$  __$$\ $$  __$$ |$$  __$$\ $$  __$$\ 
 \____$$\  \____$$\ $$  __$$ |\______|$$ |\_$$ | $$$$$$$ |$$ |  \__|$$ /  $$ |$$$$$$$$ |$$ |  $$ |
$$\   $$ |$$\   $$ |$$ |  $$ |        $$ |  $$ |$$  __$$ |$$ |      $$ |  $$ |$$   ____|$$ |  $$ |
\$$$$$$  |\$$$$$$  |$$ |  $$ |        \$$$$$$  |\$$$$$$$ |$$ |      \$$$$$$$ |\$$$$$$$\ $$ |  $$ |
 \______/  \______/ \__|  \__|         \______/  \_______|\__|       \_______| \_______|\__|  \__|`

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, title)
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return model{}, []tea.ProgramOption{tea.WithAltScreen()}
}
