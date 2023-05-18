package posts

import (
	"net/http"

	"github.com/livebud/weblog/bud/pkg/table/enum"
	"github.com/livebud/weblog/bud/pkg/table/user"
	"github.com/livebud/weblog/view"
	"github.com/matthewmueller/bud/di"
	"github.com/matthewmueller/bud/web/router"
	"github.com/xeonx/timeago"

	"github.com/gorilla/sessions"
	"github.com/livebud/bud/framework/controller/controllerrt/request"
	pogo "github.com/livebud/weblog/bud/pkg/table"
	"github.com/livebud/weblog/bud/pkg/table/post"
	"github.com/livebud/weblog/middleware/csrf"
	"github.com/livebud/weblog/view/posts"
	"gopkg.in/go-playground/validator.v9"
)

func Provide(in di.Injector) (*Controller, error) {
	db, err := di.Load[pogo.DB](in)
	if err != nil {
		return nil, err
	}
	model, err := di.Load[*post.Model](in)
	if err != nil {
		return nil, err
	}
	store, err := di.Load[sessions.Store](in)
	if err != nil {
		return nil, err
	}
	postView, err := di.Load[*posts.View](in)
	if err != nil {
		return nil, err
	}
	view, err := di.Load[view.View](in)
	if err != nil {
		return nil, err
	}
	return New(db, model, store, postView, view), nil
}

func Register(in di.Injector, router *router.Router) error {
	controller, err := di.Load[*Controller](in)
	if err != nil {
		return err
	}
	router.Get("/", controller.Index)
	router.Get("/new", controller.New)
	router.Post("/", controller.Create)
	router.Get("/:slug", controller.Show)
	router.Get("/:slug/edit", controller.Edit)
	router.Patch("/:slug", controller.Update)
	router.Delete("/:slug", controller.Delete)
	return nil
}

func New(db pogo.DB, model *post.Model, store sessions.Store, postView *posts.View, view view.View) *Controller {
	return &Controller{db, model, store, postView, view}
}

type Controller struct {
	db       pogo.DB
	model    *post.Model
	session  sessions.Store
	postView *posts.View
	view     view.View
}

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	session, err := c.session.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	loggedIn := false
	if _, ok := session.Values["user_id"].(int); ok {
		loggedIn = true
	}
	_ = loggedIn
	// ^^^ should go in middleware
	allPosts, err := c.model.FindMany(c.db, post.NewOrder().CreatedAt(post.DESC))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	props := map[string]any{
		"Title": "Posts Index",
		"Posts": allPosts,
	}
	if err := c.view.Render(r.Context(), w, "posts/index", props); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// var indexPosts []*posts.IndexPost
	// for _, post := range allPosts {
	// 	indexPosts = append(indexPosts, &posts.IndexPost{
	// 		Title:     post.Title,
	// 		Slug:      post.Slug,
	// 		CreatedAt: post.CreatedAt.Format(time.Kitchen),
	// 	})
	// }
	// c.postView.Index.Render(w, &posts.Index{
	// 	CSRF:     csrf.Token(r),
	// 	LoggedIn: loggedIn,
	// 	Posts:    indexPosts,
	// })
}

func (c *Controller) IndexJSON(w http.ResponseWriter, r *http.Request) {
	// validate the input
	// fetch a list of resources
	// render json
}

func (c *Controller) New(w http.ResponseWriter, r *http.Request) {
	session, err := c.session.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Ensure we're logged in
	_, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	c.postView.New.Render(w, &posts.New{
		CSRF: csrf.Token(r),
	})
}

var valid = validator.New()

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string          `json:"title" validate:"required"`
		Slug   string          `json:"slug" validate:"required"`
		Status enum.PostStatus `json:"status" validate:"required"`
		Body   string          `json:"body" validate:"required"`
	}
	if err := request.Unmarshal(r, &input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := valid.Struct(input); err != nil {
		// TODO: show form errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, err := c.session.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Ensure we're logged in
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	post, err := c.model.Insert(c.db, post.New().AuthorID(userID).Title(input.Title).Slug(input.Slug).Status(input.Status).Body(input.Body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/"+post.Slug, http.StatusFound)
}

func (c *Controller) CreateJSON(w http.ResponseWriter, r *http.Request) {
	// validate the input
	// create the resource
	// render json
}

func (c *Controller) Show(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Slug string `json:"slug" validate:"required"`
	}
	if err := request.Unmarshal(r, &input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := valid.Struct(input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post, err := c.model.FindBySlug(c.db, input.Slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	author, err := user.FindByID(c.db, post.AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, err := c.session.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isAuthor := false
	if userID, ok := session.Values["user_id"].(int); ok && userID == author.ID {
		isAuthor = true
	}
	if err := c.postView.Show.Render(w, &posts.Show{
		CSRF: csrf.Token(r),
		Post: &posts.ShowPost{
			Title:      post.Title,
			Slug:       post.Slug,
			Body:       post.Body,
			Status:     string(post.Status),
			Author:     author.Name,
			IsAuthor:   isAuthor,
			CreatedAgo: timeago.English.Format(post.CreatedAt),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) ShowJSON(w http.ResponseWriter, r *http.Request) {
	// validate the input
	// fetch the resource
	// render json
}

func (c *Controller) Edit(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Slug string `json:"slug" validate:"required"`
	}
	if err := request.Unmarshal(r, &input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := valid.Struct(input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post, err := c.model.FindBySlug(c.db, input.Slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	author, err := user.FindByID(c.db, post.AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, err := c.session.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isAuthor := false
	if userID, ok := session.Values["user_id"].(int); ok && userID == author.ID {
		isAuthor = true
	} else {
		http.Redirect(w, r, "/"+post.Slug, http.StatusFound)
	}
	if err := c.postView.Edit.Render(w, &posts.Edit{
		CSRF: csrf.Token(r),
		Post: &posts.EditPost{
			Title:      post.Title,
			Slug:       post.Slug,
			Body:       post.Body,
			Status:     string(post.Status),
			Author:     author.Name,
			IsAuthor:   isAuthor,
			CreatedAgo: timeago.English.Format(post.CreatedAt),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string          `json:"title" validate:"required"`
		Slug   string          `json:"slug" validate:"required"`
		Status enum.PostStatus `json:"status" validate:"required"`
		Body   string          `json:"body" validate:"required"`
	}
	if err := request.Unmarshal(r, &input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := valid.Struct(input); err != nil {
		// TODO: show form errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, err := c.session.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Ensure we're logged in
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	post, err := c.model.UpdateBySlug(c.db, input.Slug, post.New().AuthorID(userID).Title(input.Title).Slug(input.Slug).Status(input.Status).Body(input.Body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/"+post.Slug, http.StatusFound)
}

func (c *Controller) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	// validate the input
	// update the resource
	// render json
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Slug string `json:"slug" validate:"required"`
	}
	if err := request.Unmarshal(r, &input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := valid.Struct(input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, err := c.session.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Ensure we're logged in
	if _, ok := session.Values["user_id"].(int); !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	if _, err := c.model.DeleteBySlug(c.db, input.Slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *Controller) DeleteJSON(w http.ResponseWriter, r *http.Request) {
	// validate the input
	// delete the resource
	// render json
}
