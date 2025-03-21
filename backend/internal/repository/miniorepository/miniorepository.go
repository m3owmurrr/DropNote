package miniorepository

import (
	"context"
	"log"
	"os"

	"github.com/m3owmurrr/DropNote/backend/internal/model"
	"github.com/minio/minio-go/v7"
)

type NoteRepository struct {
	s3 *minio.Client
}

func NewNoteRepository(s3 *minio.Client) *NoteRepository {
	return &NoteRepository{
		s3: s3,
	}
}

func (nr *NoteRepository) CreateNote(note *model.Note) error {
	tmpFile, err := os.CreateTemp("", "note_*.txt")
	if err != nil {
		log.Println("Ошибка создания временного файла:", err)
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(note.Text)
	if err != nil {
		log.Println("Ошибка записи в файл:", err)
		return err
	}
	tmpFile.Close()

	_, err = nr.s3.FPutObject(context.TODO(), "notes-bucket", note.NoteId, tmpFile.Name(), minio.PutObjectOptions{
		ContentType: "text/plain",
	})
	if err != nil {
		log.Println("Ошибка загрузки в MinIO:", err)
		return err
	}

	return nil
}

// func (nr *NoteRepository) GetNote(id string) (*model.Note, error) {
// 	noteFile := fmt.Sprintf("%s.txt", id)

// 	err := nr.s3.FGetObject(context.TODO(), "notes-bucket", id, noteFile, minio.GetObjectOptions{})
// 	if err != nil {
// 		log.Println("Ошибка при получении объекта:", err)
// 		return nil, err
// 	}

// 	data, err := os.ReadFile(noteFile)
// 	if err != nil {
// 		log.Println("Ошибка при чтении файла:", err)
// 		return nil, err
// 	}

// 	defer os.Remove(noteFile)

// 	note := model.Note{
// 		NoteId: id,
// 		Text:   string(data),
// 	}

// 	return &note, nil
// }

// func (nr *NoteRepository) GetPublicNotes() (model.Notes, error) {
// 	return nil, nil
// }

func (nr *NoteRepository) DeleteNote(id string) error {
	err := nr.s3.RemoveObject(context.TODO(), "notes-bucket", id, minio.RemoveObjectOptions{})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
