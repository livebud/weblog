package users

import (
	_ "embed"
	"net/http"

	"github.com/livebud/weblog/internal/gohtml"
)

func Load() *View {
	return &View{
		Signup: &signupView{},
		Index:  &indexView{},
	}
}

type View struct {
	Signup *signupView
	Index  *indexView
}

//go:embed signup.gohtml
var signup string

var signupTemplate = gohtml.MustParse("signup.gohtml", signup)

type signupView struct {
}

type Signup struct {
	CSRF string
}

func (v *signupView) Render(w http.ResponseWriter, props *Signup) error {
	return signupTemplate.Render(w, props)
}

//go:embed index.gohtml
var index string

var indexTemplate = gohtml.MustParse("index.gohtml", index)

type indexView struct {
}

type IndexProps struct {
	CSRF  string
	Name  string
	Email string
}

func (v *indexView) Render(w http.ResponseWriter, props *IndexProps) error {
	return indexTemplate.Render(w, props)
}
