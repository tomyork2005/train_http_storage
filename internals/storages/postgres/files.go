package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"train_http_storage/internals/models"
)

func (s *Storage) GetAll(ctx context.Context, id int64) ([]models.File, error) {
	const op = "storage.postgres.getAll"

	files := make([]models.File, 0)
	err := s.db.SelectContext(
		ctx,
		&files,
		`SELECT * FROM files WHERE user_id = $1`,
		id)
	if err != nil {
		slog.Error("wrong with getAll", slog.Int64("id", id))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return files, nil
}
