package main

import (
	"log"
	"net/http"

	"notes-app/internal/config"
	"notes-app/internal/db"
	"notes-app/internal/handler"
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

	db.RunMigrations(cfg.DatabaseURL)

	noteRepo := repository.NewPostgresNoteRepository(pool)
	noteService := service.NewNoteService(noteRepo)
	noteHandler := handler.NewNoteHandler(noteService)

	userRepo := repository.NewPostgresUserRepository(pool)
	refreshTokenRepo := repository.NewPostgresRefreshTokenRepository(pool)
	userService := service.NewUserService(userRepo, refreshTokenRepo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userService)

	r := router.NewRouter(noteHandler, userHandler, cfg.JWTSecret)

	log.Println("Server running on :" + cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
