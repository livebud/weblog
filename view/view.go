package view

import (
	"embed"
	"os"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/view"
)

//go:embed *.gohtml **/*.gohtml
var fsys embed.FS

var pages = view.Pages{
	"posts/index": {
		Layout: `layout.gohtml`,
		Frames: []string{
			`posts/frame.gohtml`,
		},
		Path: `posts/index.gohtml`,
	},
}

type View = view.View

func Provider(in di.Injector) {
	di.Provide[view.FS](in, provideFS)
	di.Provide[view.Pages](in, providePages)
}

func provideFS(in di.Injector) (view.FS, error) {
	return os.DirFS("view"), nil
	// return fsys, nil
}

func providePages(in di.Injector) (view.Pages, error) {
	return pages, nil
}
