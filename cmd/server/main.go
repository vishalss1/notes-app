package main

import (
	"log"
	"net/http"

	"notes-app/internal/config"
	"notes-app/internal/db"
	"notes-app/internal/handler"
	"notes-app/internal/middleware"
	"notes-app/internal/repository"
	"notes-app/internal/router"
	"notes-app/internal/service"
)

func main() {
	cfg := config.Load()

	pool, err := db.NewPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	noteRepo := repository.NewPostgresNoteRepository(pool)

	noteService := service.NewNoteService(noteRepo)

	noteHandler := handler.NewNoteHandler(noteService)

	r := router.NewRouter(noteHandler)

	wrapped := middleware.Logging(r)

	log.Println("Server running on :" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, wrapped))
}
