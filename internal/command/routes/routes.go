package routes

import (
	"context"
	"fmt"

	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/web/router"
)

func Provide(in di.Injector) (*Command, error) {
	router, err := di.Load[*router.Router](in)
	if err != nil {
		return nil, err
	}
	return New(router), nil
}

func Register(in di.Injector, cli *cli.CLI) error {
	cmd, err := di.Load[*Command](in)
	if err != nil {
		return err
	}
	sub := cli.Command("routes", "list the routes")
	sub.Run(cmd.Routes)
	return nil
}

func New(router *router.Router) *Command {
	return &Command{router}
}

type Command struct {
	router *router.Router
}

func (c *Command) Routes(ctx context.Context) error {
	for _, route := range c.router.List() {
		fmt.Println(route.Method, route.Path, route.Name)
	}
	return nil
}
