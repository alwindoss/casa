package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/alwindoss/casa/internal/service"
)

type PageHandler interface {
	ShowHome(w http.ResponseWriter, r *http.Request)
}

func NewPageHandler(sess *scs.SessionManager, svc service.UserService) PageHandler {
	tc, err := createTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	ph := &pageHandler{
		sess:          sess,
		userSvc:       svc,
		templateCache: tc,
	}
	return ph
}

type pageHandler struct {
	sess *scs.SessionManager
	// userRepo repository.UserRepositoy
	userSvc       service.UserService
	templateCache map[string]*template.Template
}

func (ph *pageHandler) ShowHome(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	ph.sess.Put(r.Context(), "remote-ip", remoteIP)
	d := &TemplateData{
		Title: "Casa | Home",
	}
	ph.renderTemplate(w, r, "home.page.tmpl", d)
}
