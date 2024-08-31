package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alwindoss/casa/internal/handler"
	"github.com/alwindoss/casa/internal/repository"
	"github.com/alwindoss/casa/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() error {
	err := godotenv.Load()
	if err != nil {
		err = fmt.Errorf("error loading .env file: %w", err)
		return err
	}
	dbURL := os.Getenv("DB_URL")
	addr := os.Getenv("CASA_ADDR")

	sess := scs.New()
	sess.Lifetime = 24 * time.Hour
	sess.Cookie.Persist = true
	sess.Cookie.SameSite = http.SameSiteLaxMode
	// Secure should be set to true in production and use HTTPS
	sess.Cookie.Secure = false

	mux := chi.NewRouter()

	// A good base middleware stack
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	mux.Use(middleware.Timeout(60 * time.Second))

	// Custom Middlewares
	mux.Use(handler.WriteToConsole)
	mux.Use(handler.NoSurf)
	mux.Use(sess.LoadAndSave)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return err
	}
	userRepo := repository.NewGORMUserRepository(db)
	userSvc := service.NewUserService(userRepo)

	ph := handler.NewPageHandler(sess, userSvc)

	setupRoutes(mux, ph)

	return http.ListenAndServe(addr, mux)
}

func setupRoutes(m *chi.Mux, ph handler.PageHandler) {
	m.Get("/", ph.ShowHome)
}
