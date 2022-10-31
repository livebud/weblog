package web

import (
	"net/http"

	"github.com/livebud/bud/package/middleware"
	"github.com/livebud/bud/package/router"
	pogo "github.com/livebud/weblog/bud/pkg/table"
	"github.com/livebud/weblog/bud/pkg/table/post"
	"github.com/livebud/weblog/bud/pkg/table/user"
	"github.com/livebud/weblog/controller/posts"
	"github.com/livebud/weblog/controller/sessions"
	"github.com/livebud/weblog/controller/users"
	"github.com/livebud/weblog/env"
	"github.com/livebud/weblog/middleware/csrf"
	"github.com/livebud/weblog/middleware/httpwrap"
	"github.com/livebud/weblog/middleware/session"
	postview "github.com/livebud/weblog/view/posts"
	sessionview "github.com/livebud/weblog/view/sessions"
	userview "github.com/livebud/weblog/view/users"
)

func New(db pogo.DB, env *env.Env, r *router.Router, sessionStore session.Store) Handler {
	handler := middleware.Compose(
		httpwrap.New(),
		csrf.New(env),
		middleware.MethodOverride(),
		r,
	)

	session := session.New(env)

	// users
	users := users.New(
		db,
		&user.Model{},
		userview.Load(),
		sessionStore,
	)
	r.Get("/users", session.Middleware(http.HandlerFunc(users.Index)))
	r.Get("/signup", session.Middleware(http.HandlerFunc(users.Signup)))
	r.Post("/users", session.Middleware(http.HandlerFunc(users.Create)))

	// sessions
	sessions := sessions.New(db, &user.Model{}, sessionStore, sessionview.New())
	r.Get("/login", session.Middleware(http.HandlerFunc(sessions.New)))
	r.Post("/login", session.Middleware(http.HandlerFunc(sessions.Create)))
	r.Delete("/logout", session.Middleware(http.HandlerFunc(sessions.Delete)))

	// posts
	posts := posts.New(
		db,
		&post.Model{},
		sessionStore,
		postview.Load(),
	)
	r.Get("/", session.Middleware(http.HandlerFunc(posts.Index)))
	r.Get("/new", session.Middleware(http.HandlerFunc(posts.New)))
	r.Get("/.json", session.Middleware(http.HandlerFunc(posts.IndexJSON)))
	r.Post("/", session.Middleware(http.HandlerFunc(posts.Create)))
	r.Post("/.json", session.Middleware(http.HandlerFunc(posts.CreateJSON)))
	r.Get("/:slug", session.Middleware(http.HandlerFunc(posts.Show)))
	r.Get("/:slug.json", session.Middleware(http.HandlerFunc(posts.ShowJSON)))
	r.Get("/:slug/edit", session.Middleware(http.HandlerFunc(posts.Edit)))
	r.Patch("/:slug", session.Middleware(http.HandlerFunc(posts.Update)))
	r.Patch("/:slug.json", session.Middleware(http.HandlerFunc(posts.UpdateJSON)))
	r.Delete("/:slug", session.Middleware(http.HandlerFunc(posts.Delete)))
	r.Delete("/:slug.json", session.Middleware(http.HandlerFunc(posts.DeleteJSON)))

	// return the handler
	return handler.Middleware(http.NotFoundHandler())
}

type Handler = http.Handler

func NewServer(env *env.Env, handler Handler) *Server {
	return &http.Server{
		Addr:    env.ListenAddr,
		Handler: handler,
	}
}

type Server = http.Server
