package repository

import (
	"ms_template/internal/domain"
	"sync"
)



type Postgres struct {
	notes 		map[string]domain.Note
	mu 			sync.RWMutex
}

func NewPostgresRepo() *Postgres {
	return &Postgres{}
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

func (p *Postgres) AddNote(note domain.Note) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.notes[note.ID] = note
}