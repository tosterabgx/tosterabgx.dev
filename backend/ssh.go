package main

import (
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/logging"
)

func newSSHServer(addr string) (*ssh.Server, error) {
	return wish.NewServer(
		wish.WithAddress(addr),
		wish.WithMiddleware(
			helloMiddleware,
			logging.Middleware(),
		),
	)
}

func helloMiddleware(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		wish.Println(s, "Hello, "+s.User()+"!")
		next(s)
	}
}
