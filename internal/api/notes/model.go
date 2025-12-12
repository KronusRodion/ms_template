package notes

import (
	"context"
	"log/slog"
	"ms_template/internal/api/notes/repository"
	"ms_template/internal/api/notes/usecase"
	"ms_template/internal/domain"
)

type NoteServer struct {
	log     *slog.Logger
	usecase usecase.NoteUsecase
}

func NewServer(log *slog.Logger) *NoteServer {
	
	repo := repository.NewPostgresRepo()
	usecase := usecase.NewBasic(repo)

	return &NoteServer{usecase: usecase, log: log}
}

func (n *NoteServer) AddNote(ctx context.Context, note domain.Note) {
	n.usecase.AddNote(note)
}

func (n *NoteServer) GetNotes(ctx context.Context, userID string) []domain.Note {
	return n.usecase.GetNotes(userID)
}
