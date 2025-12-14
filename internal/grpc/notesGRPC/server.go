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
	AddNote(ctx context.Context, note domain.Note) string
	GetNotes(ctx context.Context, userID string) []domain.Note
}




func Register(grpcServer *grpc.Server, nt NoteServer) {
	notes.RegisterNotesServer(grpcServer, &ServerApi{noteServer: nt})
}


func (s *ServerApi) AddNote(ctx context.Context, in *notes.AddNoteRequest) (*notes.AddNoteResponse, error){
	note := domain.Note{
		UserID: in.UserID,
		Title: in.Note.Title,
		Content: in.Note.Content,
	}

	id := s.noteServer.AddNote(ctx, note)

	out := notes.AddNoteResponse{
		Id: id,
		Title: note.Title,
		Content: note.Content,
	}
	return &out, nil
}


func (s *ServerApi) GetNotes(ctx context.Context, in *notes.GetNotesRequest) (*notes.GetNotesResponse, error){
	
	
	noteArr := s.noteServer.GetNotes(ctx, in.UserID)

	result := make([]*notes.Note, len(noteArr))

	for i,v := range noteArr {
		result[i] = &notes.Note{
			Id: v.ID,
			Title: v.Title,
			Content: v.Content,
		}
	}

	out := notes.GetNotesResponse{
		Notes: result,
	}

	return &out, nil
}

