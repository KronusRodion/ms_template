package usecase

import "ms_template/internal/domain"

type NoteUsecase interface {
	AddNote(domain.Note)
	GetNotes(userID string) []domain.Note
}