package gohtml

import (
	"html/template"
	"io"
)

func MustParse(name, code string) *Template {
	tpl, err := Parse(name, code)
	if err != nil {
		panic(err)
	}
	return tpl
}

// Parse parses Go code
func Parse(name, code string) (*Template, error) {
	tpl, err := template.New(name).Parse(code)
	if err != nil {
		return nil, err
	}
	return &Template{tpl}, nil
}

// Template struct
type Template struct {
	tpl *template.Template
}

// Render the template
func (t *Template) Render(w io.Writer, props interface{}) error {
	return t.tpl.Execute(w, props)
}
