package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"notes-app/internal/handler"
)

func NewRouter(noteHandler *handler.NoteHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API routes
	r.Route("/notes", func(r chi.Router) {
		r.Get("/", noteHandler.GetAll)
		r.Post("/", noteHandler.Create)
		r.Get("/{id}", noteHandler.GetByID)
		r.Put("/{id}", noteHandler.Update)
		r.Delete("/{id}", noteHandler.Delete)
	})

	// Serve static frontend from ./web
	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fileServer)

	return r
}
