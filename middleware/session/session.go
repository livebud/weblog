package session

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/livebud/weblog/env"
)

type Store = sessions.Store

func New(env *env.Env) Middleware {
	store := sessions.NewCookieStore([]byte(env.SessionKey))
	return Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			// Save the session
			session, err := store.Get(r, "user")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Update the visitor count
			if n, ok := session.Values["visits"].(int); ok {
				session.Values["visits"] = n + 1
			} else {
				session.Values["visits"] = 1
			}
			session.Save(r, w)
		})
	})
}

type Middleware func(next http.Handler) http.Handler

func (m Middleware) Middleware(next http.Handler) http.Handler {
	return m(next)
}
