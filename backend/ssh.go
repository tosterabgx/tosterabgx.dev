package main

import (
	"database/sql"
	"net"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	bm "charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
)

func newSSHServer(addr string, db *sql.DB) (*ssh.Server, error) {
	handler := func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		host, _, err := net.SplitHostPort(s.RemoteAddr().String())
		if err != nil {
			host = s.RemoteAddr().String()
		}
		return newModel(db, host), nil
	}

	return wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath("ssh_host_key"),
		wish.WithIdleTimeout(2*time.Minute),
		wish.WithMaxTimeout(10*time.Minute),
		wish.WithMiddleware(
			bm.Middleware(handler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
}
