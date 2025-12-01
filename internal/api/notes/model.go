package notes

import (
	"context"
	"log/slog"
	"ms_template/internal/api/notes/usecase"
	"ms_template/internal/domain"
)

type NoteServer struct {
	log     *slog.Logger
	usecase usecase.NoteUsecase
}

func NewServer() *NoteServer {
	// Здесь создаем репо и usecase и передаем в NoteServer
	// repository := repository.NoteRepository
	// usecase := usecase.NoteUsecase

	return &NoteServer{}
}

func (n *NoteServer) AddNote(ctx context.Context, note domain.Note) {
	n.usecase.AddNote(note)
}

func (n *NoteServer) GetNotes(ctx context.Context, userID string) []domain.Note {
	return n.usecase.GetNotes(userID)
}
