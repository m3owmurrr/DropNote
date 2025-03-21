package service

import (
	"github.com/m3owmurrr/DropNote/backend/internal/model"
	"github.com/m3owmurrr/DropNote/backend/internal/repository"
)

type NoteService struct {
	metaDB    repository.Repository
	s3Storage repository.Repository
}

func NewNoteSevice(metaDB repository.Repository, s3Storage repository.Repository) *NoteService {
	return &NoteService{
		metaDB:    metaDB,
		s3Storage: s3Storage,
	}
}

func (ns *NoteService) CreateNote(note *model.Note) error {
	if err := ns.s3Storage.CreateNote(note); err != nil {
		return err
	}

	if err := ns.metaDB.CreateNote(note); err != nil {
		ns.s3Storage.DeleteNote(note.NoteId)
		return err
	}
	return nil
}

// func (ns *NoteService) GetNote(id string) (*model.Note, error) {
// 	note, err := ns.s3Storage.GetNote(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return note, nil
// }

// func (ns *NoteService) GetPublicNotes() (model.Notes, error) {
// 	metaNotes, err := ns.metaDB.GetPublicNotes()
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Println("Fetched meta notes:", metaNotes)

// 	var notes model.Notes
// 	for _, metaInfo := range metaNotes {
// 		note, err := ns.s3Storage.GetNote(metaInfo.NoteId)
// 		if err != nil {
// 			continue
// 		}

// 		notes = append(notes, note)
// 	}

// 	return notes, nil
// }
