package seed

import (
	"context"
	"fmt"

	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/db"
	"github.com/matthewmueller/bud/di"
)

func Provide(in di.Injector) (*Command, error) {
	db, err := di.Load[db.DB](in)
	if err != nil {
		return nil, err
	}
	return New(db), nil
}

func Register(in di.Injector, cli *cli.CLI) error {
	cmd, err := di.Load[*Command](in)
	if err != nil {
		return err
	}
	sub := cli.Command("seed", "seed the database")
	sub.Run(cmd.Seed)
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
