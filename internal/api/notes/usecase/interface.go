package usecase

import "ms_template/internal/domain"

type NoteUsecase interface {
	AddNote(domain.Note) string
	GetNotes(userID string) []domain.Note
}