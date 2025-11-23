package service

import (
	"context"
	"gRPC-server/internal/model"
	"gRPC-server/internal/repository"
	notepb "gRPC-server/proto"
	"log"
)

type NoteService struct {
	Repo repository.Repository
	notepb.UnimplementedNoteServiceServer
}

func NewService(repository2 repository.Repository) *NoteService {
	return &NoteService{Repo: repository2}
}

func (noteService *NoteService) CreateNote(ctx context.Context, req *notepb.CreateNoteRequest) (*notepb.CreateNoteResponse, error) {
	id, err := noteService.Repo.CreateNote(ctx, model.Note{UserID: req.UserId, Text: req.NoteText})
	if err != nil {
		return nil, err
	}

	log.Printf("Note successfully create by userID=%v", req.UserId)

	return &notepb.CreateNoteResponse{NoteId: id}, nil
}
