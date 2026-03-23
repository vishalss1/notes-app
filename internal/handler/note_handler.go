package handler

import (
	"encoding/json"
	"net/http"

	"notes-app/internal/middleware"
	"notes-app/internal/model"
	"notes-app/internal/repository"
	"notes-app/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NoteHandler struct {
	service *service.NoteService
}

func NewNoteHandler(s *service.NoteService) *NoteHandler {
	return &NoteHandler{service: s}
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	note := model.Note{
		ID:      uuid.New().String(),
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	created, err := h.service.Create(r.Context(), note)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	notes, err := h.service.GetAll(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch notes"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	note, err := h.service.GetByID(r.Context(), id, userID)
	if err != nil {
		if err == repository.ErrNoteNotFound {
			http.Error(w, `{"error":"note not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"failed to fetch note"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	note := model.Note{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
	}

	updated, err := h.service.Update(r.Context(), note, userID)
	if err != nil {
		if err == repository.ErrNoteNotFound {
			http.Error(w, `{"error":"note not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	err := h.service.Delete(r.Context(), id, userID)
	if err != nil {
		if err == repository.ErrNoteNotFound {
			http.Error(w, `{"error":"note not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"delete failed"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
