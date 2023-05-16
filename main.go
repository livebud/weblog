package main

import (
	"os"

	"github.com/livebud/weblog/internal/injector"
	"github.com/matthewmueller/bud/app"
)

func main() {
	os.Exit(app.Run(injector.New()))
}
