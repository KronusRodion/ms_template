package notesGRPC

import (
	"context"
	"ms_template/gen/go/notes"
	"ms_template/internal/domain"

	"google.golang.org/grpc"
)







type ServerApi struct {
	notes.UnimplementedNotesServer
	noteServer NoteServer
}

type NoteServer interface {
	AddNote(ctx context.Context, note domain.Note)
	GetNotes(ctx context.Context, userID string) []domain.Note
}




func Register(grpcServer *grpc.Server, nt NoteServer) {
	notes.RegisterNotesServer(grpcServer, &ServerApi{noteServer: nt})
}


func (s *ServerApi) AddNote(ctx context.Context, in *notes.AddNoteRequest) (*notes.AddNoteResponse, error){
	note := domain.Note{
		Title: in.Note.Title,
		Content: in.Note.Content,
	}

	s.noteServer.AddNote(ctx, note)

	out := notes.AddNoteResponse{
		Title: note.Title,
		Content: note.Content,
	}
	return &out, nil
}


func (s *ServerApi) GetNotes(ctx context.Context, in *notes.GetNotesRequest) (*notes.GetNotesResponse, error){
	return nil, nil
}

