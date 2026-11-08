package main

import (
	"database/sql"

	tea "charm.land/bubbletea/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	bm "charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
)

func newSSHServer(addr string, db *sql.DB) (*ssh.Server, error) {
	handler := func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		return newModel(db, s.RemoteAddr().String()), nil
	}

	return wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath("ssh_host_key"),
		wish.WithMiddleware(
			bm.Middleware(handler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
}
