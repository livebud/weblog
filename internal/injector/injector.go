package injector

import (
	"errors"
	"os"

	"github.com/gorilla/sessions"
	pogo "github.com/livebud/weblog/bud/pkg/table"
	"github.com/livebud/weblog/bud/pkg/table/post"
	"github.com/livebud/weblog/controller/posts"
	"github.com/livebud/weblog/internal/command/routes"
	"github.com/livebud/weblog/internal/command/seed"
	postsview "github.com/livebud/weblog/view/posts"
	"github.com/matthewmueller/bud/cli"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/injector"
	"github.com/matthewmueller/bud/web/router"
)

func New() di.Injector {
	in := injector.New()
	pogo.Provider(in)
	post.Provider(in)
	di.Provide[*seed.Command](in, seed.Provide)
	di.Register[*cli.CLI](in, seed.Register)
	di.Provide[*routes.Command](in, routes.Provide)
	di.Register[*cli.CLI](in, routes.Register)
	di.Provide[*posts.Controller](in, posts.Provide)
	di.Register[*router.Router](in, posts.Register)
	di.Provide[sessions.Store](in, provideSessions)
	di.Provide[*postsview.View](in, postsview.Provide)
	return in
}

func provideSessions(in di.Injector) (sessions.Store, error) {
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		return nil, errors.New("SESSION_KEY is required")
	}
	return sessions.NewCookieStore([]byte(sessionKey)), nil
}
