package usecase

import (
	"ms_template/internal/api/notes/repository"
	"ms_template/internal/domain"
)


type Basic struct {
	repo 	repository.NoteRepository
}

func NewBasic(repo repository.NoteRepository) *Basic {
	return &Basic{repo: repo}
}

func (b *Basic) GetNotes(userID string) []domain.Note {
	return b.repo.GetNotes()
}

func (b *Basic) AddNote(note domain.Note) {
	b.repo.AddNote(note)
}