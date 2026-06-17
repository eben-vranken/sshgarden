package main

import (
	"errors"
	"net"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
