package router

import (
	"net/http"
	"strings"

	"notes-app/internal/handler"
)

func NewRouter(noteHandler *handler.NoteHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			noteHandler.GetAll(w, r)
		case http.MethodPost:
			noteHandler.Create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/notes/") {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			noteHandler.GetByID(w, r)
		case http.MethodPut:
			noteHandler.Update(w, r)
		case http.MethodDelete:
			noteHandler.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	return mux
}
