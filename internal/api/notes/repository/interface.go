package repository

import "ms_template/internal/domain"

type NoteRepository interface {
	AddNote(domain.Note)
	GetNotes() []domain.Note
}
