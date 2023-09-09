package posts_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/livebud/weblog/internal/injector"
	"github.com/matryer/is"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/view"
)

func TestPosts(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	in := injector.New()
	viewer, err := di.Load[*view.Viewer](in)
	is.NoErr(err)
	w := new(bytes.Buffer)
	fmt.Println(ctx, viewer, w)
	// viewer.Render(ctx, w, "posts/index")
	// posts, err := di.Load[*posts.Controller](in)

}
