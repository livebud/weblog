package main

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"

	"github.com/jackc/pgx"

	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/weblog/env"
	"github.com/livebud/weblog/web"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	env, err := env.Load()
	if err != nil {
		return fmt.Errorf("unable to load environment. %s", err)
	}
	pgconfig, err := pgx.ParseURI(env.DatabaseURL)
	if err != nil {
		return err
	}
	db, err := pgx.Connect(pgconfig)
	if err != nil {
		return err
	}
	defer db.Close()
	cookieStore := sessions.NewCookieStore([]byte(env.SessionKey))
	handler := web.New(db, env, router.New(), cookieStore)
	server := web.NewServer(env, handler)
	console.Info("Listening on http://localhost" + env.ListenAddr)
	return server.ListenAndServe()
}
