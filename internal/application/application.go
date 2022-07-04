package application

import (
	"github.com/storskegg/r53transfer/internal/clients"
	"github.com/urfave/cli/v2"
)

type Application interface {
	Run([]string) error
}

type app struct {
	App     cli.App
	Clients clients.Clients
}

func (a *app) Run(args []string) (err error) {
	return
}

func New() Application {
	a := &app{}
	a.App = cli.App{}

	return a
}
