package users

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/livebud/bud/framework/controller/controllerrt/request"
	pogo "github.com/livebud/weblog/bud/pkg/table"
	"github.com/livebud/weblog/bud/pkg/table/user"
	"github.com/livebud/weblog/middleware/csrf"
	"github.com/livebud/weblog/view/users"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

func New(db pogo.DB, model *user.Model, view *users.View, sessions sessions.Store) *Controller {
	return &Controller{db, model, view, sessions}
}

type Controller struct {
	db       pogo.DB
	model    *user.Model
	view     *users.View
	sessions sessions.Store
}

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	session, err := c.sessions.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Missing
	userId, ok := session.Values["user_id"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Invalid type
	userIdInt, ok := userId.(int)
	if !ok {
		http.Error(w, "invalid user id", http.StatusInternalServerError)
		return
	}
	usr, err := c.model.FindByID(c.db, userIdInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := c.view.Index.Render(w, &users.IndexProps{
		CSRF:  csrf.Token(r),
		Name:  usr.Name,
		Email: usr.Email,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := c.view.Signup.Render(w, &users.Signup{
		CSRF: csrf.Token(r),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var valid = validator.New()

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
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
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	usr, err := c.model.Insert(c.db, user.New().Name(input.Name).Email(input.Email).Password(string(password)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Save the session
	session, err := c.sessions.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user_id"] = usr.ID
	// Redirect to the user's profile
	http.Redirect(w, r, "/", http.StatusFound)
}
