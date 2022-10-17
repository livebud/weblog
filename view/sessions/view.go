package sessions

import (
	_ "embed"
	"net/http"

	"github.com/livebud/weblog/internal/gohtml"
)

func New() *View {
	return &View{
		Login: &loginView{},
	}
}

type View struct {
	Login *loginView
}

//go:embed login.gohtml
var login string

var loginTemplate = gohtml.MustParse("login.gohtml", login)

type loginView struct {
}

type Login struct {
	CSRF string
}

func (v *loginView) Render(w http.ResponseWriter, props *Login) error {
	return loginTemplate.Render(w, props)
}
