package posts

import (
	_ "embed"
	"net/http"

	"github.com/livebud/weblog/internal/gohtml"
)

func Load() *View {
	return &View{
		Index: &indexView{},
		Show:  &showView{},
		New:   &newView{},
		Edit:  &editView{},
	}
}

type View struct {
	Index *indexView
	Show  *showView
	New   *newView
	Edit  *editView
}

//go:embed index.gohtml
var index string

var indexTemplate = gohtml.MustParse("index.gohtml", index)

type indexView struct {
}

type Index struct {
	CSRF     string
	LoggedIn bool
	Posts    []*IndexPost
}

type IndexPost struct {
	Title     string
	Slug      string
	Author    string
	CreatedAt string
}

func (i *indexView) Render(w http.ResponseWriter, props *Index) error {
	return indexTemplate.Render(w, props)
}

//go:embed show.gohtml
var show string

var showTemplate = gohtml.MustParse("show.gohtml", show)

type showView struct {
}

type Show struct {
	CSRF string
	Post *ShowPost
}

type ShowPost struct {
	Title      string
	Slug       string
	Author     string
	Status     string
	IsAuthor   bool
	Body       string
	CreatedAgo string
}

func (v *showView) Render(w http.ResponseWriter, props *Show) error {
	return showTemplate.Render(w, props)
}

//go:embed new.gohtml
var new string

var newTemplate = gohtml.MustParse("new.gohtml", new)

type newView struct {
}

type New struct {
	CSRF string
}

func (v *newView) Render(w http.ResponseWriter, props *New) error {
	return newTemplate.Render(w, props)
}

//go:embed edit.gohtml
var edit string

var editTemplate = gohtml.MustParse("edit.gohtml", edit)

type editView struct {
}

type Edit struct {
	CSRF string
	Post *EditPost
}

type EditPost struct {
	Title      string
	Slug       string
	Author     string
	Status     string
	IsAuthor   bool
	Body       string
	CreatedAgo string
}

func (v *editView) Render(w http.ResponseWriter, props *Edit) error {
	return editTemplate.Render(w, props)
}
