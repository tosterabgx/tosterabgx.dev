package main

import (
	"database/sql"
	"log"
	"net"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	bm "charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
)

func newSSHServer(addr string, db *sql.DB, adminKey ssh.PublicKey) (*ssh.Server, error) {
	handler := func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		host, _, err := net.SplitHostPort(s.RemoteAddr().String())
		if err != nil {
			host = s.RemoteAddr().String()
		}

		if isAdmin(s, adminKey) {
			log.Printf("admin connected from %s", host)

		}
		return newModel(db, host), nil
	}

	return wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath("ssh_host_key"),
		wish.WithIdleTimeout(2*time.Minute),
		wish.WithMaxTimeout(10*time.Minute),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			bm.Middleware(handler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
}

func isAdmin(s ssh.Session, adminKey ssh.PublicKey) bool {
	return adminKey != nil &&
		s.PublicKey != nil &&
		ssh.KeysEqual(s.PublicKey(), adminKey)
}
