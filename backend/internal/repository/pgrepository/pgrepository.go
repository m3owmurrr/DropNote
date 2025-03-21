package pgrepository

import (
	"database/sql"
	"errors"

	"github.com/m3owmurrr/DropNote/backend/internal/model"
)

type NoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{
		db: db,
	}
}

func (nr *NoteRepository) CreateNote(note *model.Note) error {
	insert := "INSERT INTO Notes (note_id, public) VALUES ($1, $2)"

	_, err := nr.db.Exec(insert, note.NoteId, note.IsPublic)
	if err != nil {
		return err
	}

	return nil
}

func (nr *NoteRepository) GetNote(id string) (*model.Note, error) {
	select_str := "SELECT note_id FROM Notes WHERE note_id=$1"
	row := nr.db.QueryRow(select_str, id)

	var note model.Note
	if err := row.Scan(&note.NoteId); err != nil {
		return nil, err
	}

	return &note, nil
}

func (nr *NoteRepository) GetPublicNotes() (model.Notes, error) {
	rows, err := nr.db.Query("SELECT note_id FROM Notes WHERE public = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes model.Notes
	for rows.Next() {
		var note model.Note
		if err := rows.Scan(&note.NoteId); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if notes == nil {
		return nil, errors.New("metaNotes is nil")
	}

	return notes, nil
}

func (nr *NoteRepository) DeleteNote(id string) error {
	return nil
}
