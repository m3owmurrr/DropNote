package service

import "github.com/m3owmurrr/DropNote/backend/internal/model"

type Service interface {
	CreateNote(note *model.Note) error
	// GetNote(id string) (*model.Note, error)
	// GetPublicNotes() (model.Notes, error)
}
