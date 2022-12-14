package posts

import (
	"net/http"
	"time"

	"github.com/livebud/weblog/bud/pkg/table/enum"
	"github.com/livebud/weblog/bud/pkg/table/user"
	"github.com/xeonx/timeago"

	"github.com/gorilla/sessions"
	"github.com/livebud/bud/framework/controller/controllerrt/request"
	pogo "github.com/livebud/weblog/bud/pkg/table"
	"github.com/livebud/weblog/bud/pkg/table/post"
	"github.com/livebud/weblog/middleware/csrf"
	"github.com/livebud/weblog/view/posts"
	"gopkg.in/go-playground/validator.v9"
)

func New(db pogo.DB, model *post.Model, store sessions.Store, view *posts.View) *Controller {
	return &Controller{db, model, store, view}
}

type Controller struct {
	db      pogo.DB
	model   *post.Model
	session sessions.Store
	view    *posts.View
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
	allPosts, err := c.model.FindMany(c.db, post.NewOrder().CreatedAt(post.DESC))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var indexPosts []*posts.IndexPost
	for _, post := range allPosts {
		indexPosts = append(indexPosts, &posts.IndexPost{
			Title:     post.Title,
			Slug:      post.Slug,
			CreatedAt: post.CreatedAt.Format(time.Kitchen),
		})
	}
	c.view.Index.Render(w, &posts.Index{
		CSRF:     csrf.Token(r),
		LoggedIn: loggedIn,
		Posts:    indexPosts,
	})
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
	c.view.New.Render(w, &posts.New{
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
	if err := c.view.Show.Render(w, &posts.Show{
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
	if err := c.view.Edit.Render(w, &posts.Edit{
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
