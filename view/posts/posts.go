package posts

import (
	"embed"

	"github.com/livebud/weblog/bud/pkg/table/post"
	"github.com/livebud/weblog/bud/pkg/table/user"

	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/view"
)

//go:embed *.gohtml
var fsys embed.FS

func Register(in di.Injector, viewer *view.Viewer) error {
	return viewer.Mount(fsys)
	// viewer.Add(&view.Page{
	// 	Path: "edit.gohtml",
	// 	Frames: []string{
	// 		"posts/frame.gohtml",
	// 	},
	// 	Layout: "layout.gohtml",
	// })
	// viewer.Add(&view.Page{
	// 	Path: "index.gohtml",
	// 	Frames: []string{
	// 		"posts/frame.gohtml",
	// 	},
	// 	Layout: "layout.gohtml",
	// })
	// viewer.Add(&view.Page{
	// 	Path: "new.gohtml",
	// 	Frames: []string{
	// 		"posts/frame.gohtml",
	// 	},
	// 	Layout: "layout.gohtml",
	// })
	// viewer.Add(&view.Page{
	// 	Path: "show.gohtml",
	// 	Frames: []string{
	// 		"posts/frame.gohtml",
	// 	},
	// 	Layout: "layout.gohtml",
	// })
	// return nil
}

type Edit struct {
	CSRF string
	Post *post.Post
}

type Frame struct {
}

type Index struct {
	// TODO: CSRF & User should go in middleware
	CSRF  string
	User  *user.User
	Posts []*post.Post
}

type New struct {
	CSRF string
}

type Show struct {
	CSRF string
	Post *post.Post
}
