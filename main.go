package main

import (
	"os"

	"github.com/livebud/weblog/internal/di"
	"github.com/matthewmueller/bud/app"
)

func main() {
	os.Exit(app.Run(di.New()))
}
