package csrf

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/livebud/weblog/env"
)

func New(env *env.Env) Middleware {
	return csrf.Protect([]byte(env.CSRFToken), csrf.Path("/"))
}

type Middleware func(next http.Handler) http.Handler

func (m Middleware) Middleware(next http.Handler) http.Handler {
	return m(next)
}

func Token(r *http.Request) string {
	return csrf.Token(r)
}
