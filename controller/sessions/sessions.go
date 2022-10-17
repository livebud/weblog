package sessions

import (
	"net/http"

	pogo "github.com/livebud/weblog/bud/pkg/table"

	"github.com/livebud/bud/framework/controller/controllerrt/request"
	"github.com/livebud/weblog/bud/pkg/table/user"
	"github.com/livebud/weblog/middleware/csrf"
	"github.com/livebud/weblog/middleware/session"
	"github.com/livebud/weblog/view/sessions"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

func New(db pogo.DB, model *user.Model, sessions session.Store, view *sessions.View) *Controller {
	return &Controller{db, model, sessions, view}
}

type Controller struct {
	db       pogo.DB
	model    *user.Model
	sessions session.Store
	view     *sessions.View
}

// login page
func (c *Controller) New(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := c.view.Login.Render(w, &sessions.Login{
		CSRF: csrf.Token(r),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var valid = validator.New()

// login
func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
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
	usr, err := c.model.FindByEmail(c.db, input.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(input.Password)); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

// logout
func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	session, err := c.sessions.Get(r, "user")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	delete(session.Values, "user_id")
	// Redirect to the user's profile
	http.Redirect(w, r, "/login", http.StatusFound)
}
