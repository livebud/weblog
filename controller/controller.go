package controller

import (
	"github.com/livebud/weblog/view"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/web"
	"github.com/matthewmueller/bud/web/router"
)

func Provider(in di.Injector) (*Controller, error) {
	return &Controller{}, nil
}

func Register(in di.Injector, r *router.Router) error {
	controller, err := di.Load[*Controller](in)
	if err != nil {
		return err
	}
	_ = controller
	// r.Layout(controller.Layout)
	return nil
}

type Controller struct {
}

func (c *Controller) Layout(r *web.Request[any], w web.Response[*view.Layout]) {
	w.Render(&view.Layout{
		Title: "Posts",
	})
}
