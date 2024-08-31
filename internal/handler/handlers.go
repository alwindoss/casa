package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alwindoss/casa/internal/service"
)

type PageHandler interface {
	ShowHome(w http.ResponseWriter, r *http.Request)
}

func NewPageHandler(sess *scs.SessionManager, svc service.UserService) PageHandler {
	ph := &pageHandler{
		sess:    sess,
		userSvc: svc,
	}
	return ph
}

type pageHandler struct {
	sess *scs.SessionManager
	// userRepo repository.UserRepositoy
	userSvc service.UserService
}

func (ph pageHandler) ShowHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home page: %s", time.Now())
}
