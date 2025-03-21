package model

type Note struct {
	NoteId   string `json:"id"`
	Text     string `json:"text"`
	IsPublic bool   `json:"public"`
}

type Notes []*Note
