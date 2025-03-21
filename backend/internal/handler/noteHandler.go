package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/m3owmurrr/DropNote/backend/internal/model"
	"github.com/m3owmurrr/DropNote/backend/internal/service"
)

type NoteHandler struct {
	serv   service.Service
	logger *slog.Logger
}

func NewNoteHandler(serv service.Service, logger *slog.Logger) *NoteHandler {
	return &NoteHandler{
		serv:   serv,
		logger: logger,
	}
}

// CreateNote создаёт новую заметку
//
// @Summary Создание заметки
// @Description Принимает JSON с текстом заметки, сохраняет мета-данные в PostgreSQL, а содержимое — в MinIO.
// @Tags notes
// @Accept json
// @Produce json
// @Param note body model.Note true "Данные заметки"
// @Success 201 {object} map[string]string "Заметка успешно создана, возвращает note_id"
// @Failure 400 {string} string "Некорректный JSON"
// @Failure 415 {string} string "Content-Type должен быть application/json"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/notes [post]
func (nh *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		nh.logger.Error("Content-Type is not application/json")
		http.Error(w, "Content-Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var note model.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	note.NoteId = uuid.New().String()

	if err := nh.serv.CreateNote(&note); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"note_id": note.NoteId})
}

// GetNote получает заметку по ID
//
// @Summary Получение заметки
// @Description Возвращает заметку по её идентификатору, переданному в URL.
// @Tags notes
// @Accept json
// @Produce json
// @Param note_id path string true "ID заметки"
// @Success 200 {object} model.Note "Заметка успешно получена"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Заметка не найдена"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/notes/{note_id} [get]
func (nh *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	noteID, exists := vars["note_id"]
	if !exists {
		http.Error(w, "note_id is required", http.StatusBadRequest)
		return
	}

	note, err := nh.serv.GetNote(noteID)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// func (nh *NoteHandler) GetPublicNotes(w http.ResponseWriter, r *http.Request) {
// 	notes, err := nh.serv.GetPublicNotes()
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(notes)
// }
