package posts_test

import (
	"fmt"
	"testing"

	"github.com/livebud/weblog/bud/pkg/table/post"
	"github.com/livebud/weblog/controller/posts"
	"github.com/livebud/weblog/internal/injector"
	"github.com/matthewmueller/bud/di"
)

func provideModel(in di.Injector) (*post.Model, error) {
	fmt.Println("custom model!")
	return &post.Model{}, nil
}

func TestPosts(t *testing.T) {
	in := injector.New()
	di.Provide[*post.Model](in, provideModel)
	posts, err := di.Load[*posts.Controller](in)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(posts)
}
