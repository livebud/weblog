//go:build tools
// +build tools

package main

import (
	_ "github.com/matthewmueller/migrate/cmd/migrate"
	_ "github.com/matthewmueller/pogo/cmd/pogo"
)
