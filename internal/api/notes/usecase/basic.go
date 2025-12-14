package usecase

import (
	"ms_template/internal/api/notes/repository"
	"ms_template/internal/domain"
	"time"

	"github.com/google/uuid"
)


type Basic struct {
	repo 	repository.NoteRepository
}

var _ NoteUsecase = &Basic{}

func NewBasic(repo repository.NoteRepository) *Basic {
	return &Basic{repo: repo}
}

func (b *Basic) GetNotes(userID string) []domain.Note {
	return b.repo.GetNotes()
}

func (b *Basic) AddNote(note domain.Note) string {
	note.ID = uuid.New().String()
	note.CreatedAt = time.Now()
	return b.repo.AddNote(note)
}