package repository

import (
	"ms_template/internal/domain"
	"sync"

	"github.com/google/uuid"
)

type Postgres struct {
	notes map[string]domain.Note
	mu    *sync.RWMutex
}

var _ NoteRepository = &Postgres{}

func NewPostgresRepo() *Postgres {
	mu := sync.RWMutex{}
	notes := make(map[string]domain.Note)
	return &Postgres{
		mu:    &mu,
		notes: notes,
	}
}

func (p *Postgres) GetNotes() []domain.Note {
	p.mu.RLock()
	defer p.mu.RUnlock()

	notes := make([]domain.Note, 0, len(p.notes))

	for _, note := range p.notes {
		notes = append(notes, note)
	}

	return notes
}

func (p *Postgres) AddNote(note domain.Note) string {
	if note.ID == "" {
		note.ID = uuid.New().String()
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	
	
	p.notes[note.ID] = note

	return note.ID
}
