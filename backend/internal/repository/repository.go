package repository

import "github.com/m3owmurrr/DropNote/backend/internal/model"

type Repository interface {
	CreateNote(note *model.Note) error
	GetNote(id string) (*model.Note, error)
	// GetPublicNotes() (model.Notes, error)
	DeleteNote(id string) error
}
