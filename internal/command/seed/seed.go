package seed

import (
	"context"
	"fmt"

	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/db"
	"github.com/matthewmueller/bud/di"
)

func Register(in di.Injector, cli cli.Command) error {
	db, err := di.Load[db.DB](in)
	if err != nil {
		return err
	}
	cmd := New(db)
	cli = cli.Command("seed", "seed the database")
	cli.Run(cmd.Seed)
	return nil
}

func New(db db.DB) *Command {
	return &Command{db}
}

type Command struct {
	db db.DB
}

func (c *Command) Seed(ctx context.Context) error {
	fmt.Println("seeding!", c.db)
	return nil
}
