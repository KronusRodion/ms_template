package repository

import "ms_template/internal/domain"

type NoteRepository interface {
	AddNote(domain.Note) string
	GetNotes() []domain.Note
}
