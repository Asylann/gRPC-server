package repository

import (
	"context"
	"github.com/Asylann/gRPC-server/internal/model"
)

func (repo *Repository) CreateNote(ctx context.Context, note model.Note) (string, error) {
	var id string
	err := repo.Pool.QueryRow(ctx, "INSERT INTO notes(user_id, text) VALUES($1,$2) RETURNING id", note.UserID, note.Text).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, err
}
