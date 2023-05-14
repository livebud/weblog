package routes

import (
	"context"
	"fmt"

	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/web/router"
)

func Register(in di.Injector, cli cli.Command) error {
	// db, err := di.Load[db.DB](in)
	// if err != nil {
	// 	return err
	// }
	// _ = db
	cmd := New(nil)
	cli = cli.Command("routes", "list the routes")
	cli.Run(cmd.Routes)
	return nil
}

func New(router *router.Router) *Command {
	return &Command{router}
}

type Command struct {
	router *router.Router
}

func (c *Command) Routes(ctx context.Context) error {
	fmt.Println("routes!")
	return nil
}
