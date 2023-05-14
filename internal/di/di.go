package di

import (
	"github.com/livebud/weblog/internal/command/routes"
	"github.com/livebud/weblog/internal/command/seed"
	"github.com/matthewmueller/bud/app"
	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/di"
)

func New() di.Injector {
	in := di.New()
	app.Provider(in)
	di.Register[cli.Command](in, seed.Register)
	di.Register[cli.Command](in, routes.Register)
	return in
}
