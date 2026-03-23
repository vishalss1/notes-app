package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	appHandler "notes-app/internal/handler"
	appMiddleware "notes-app/internal/middleware"
)

func NewRouter(
	noteHandler *appHandler.NoteHandler,
	userHandler *appHandler.UserHandler,
	jwtSecret string,
) *chi.Mux {

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(appMiddleware.Logging)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", userHandler.Signup)
		r.Post("/login", userHandler.Login)
		r.Post("/refresh", userHandler.Refresh)
		r.Post("/logout", userHandler.Logout)
	})

	r.Route("/notes", func(r chi.Router) {
		r.Use(appMiddleware.JWTAuth(jwtSecret))

		r.Get("/", noteHandler.GetAll)
		r.Post("/", noteHandler.Create)
		r.Get("/{id}", noteHandler.GetByID)
		r.Put("/{id}", noteHandler.Update)
		r.Delete("/{id}", noteHandler.Delete)
	})

	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fileServer)

	return r
}
