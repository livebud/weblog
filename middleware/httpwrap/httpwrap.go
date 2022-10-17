package httpwrap

import (
	"bytes"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/livebud/bud/package/middleware"
)

func New() Middleware {
	return middleware.Function(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := Wrap(w)
			next.ServeHTTP(rw, r)
			rw.Flush()
		})
	})
}

type Middleware = middleware.Middleware

// Wrap the response writer
func Wrap(w http.ResponseWriter) *ResponseWriter {
	state := new(state)
	responseWriter := httpsnoop.Wrap(w, httpsnoop.Hooks{
		WriteHeader: func(_ httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(code int) {
				state.Code = code
			}
		},
		Write: func(_ httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(p []byte) (int, error) {
				return state.Body.Write(p)
			}
		},
	})
	return &ResponseWriter{
		responseWriter,
		state,
		w,
	}
}

// Unwrap the response writer
func Unwrap(w http.ResponseWriter) (rw *ResponseWriter, ok bool) {
	rw, ok = w.(*ResponseWriter)
	return rw, ok
}

// state struct
type state struct {
	Code int
	Body bytes.Buffer
}

// ResponseWriter struct
type ResponseWriter struct {
	http.ResponseWriter

	state    *state
	original http.ResponseWriter
}

// Flush the response
func (rw *ResponseWriter) Flush() (int, error) {
	if rw.state.Code == 0 {
		rw.state.Code = http.StatusOK
	}
	rw.original.WriteHeader(rw.state.Code)
	return rw.original.Write(rw.state.Body.Bytes())
}
