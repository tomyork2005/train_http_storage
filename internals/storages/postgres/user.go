package postgres

import (
	"context"
	"fmt"
	"train_http_storage/internals/models"
)

func (s *Storage) Create(ctx context.Context, user *models.User) error {
	const op = "storage.postgres.User.Create"

	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO users (name, email, password, register_at) VALUES ($1, $2, $3, $4)`,
		user.Name,
		user.Email,
		user.Password,
		user.RegisteredAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetByCredentials(ctx context.Context, email, password string) (*models.User, error) {
	const op = "storage.postgres.User.GetByCredentials"

	var user models.User
	err := s.db.GetContext(
		ctx,
		&user,
		`SELECT * FROM users WHERE email = $1 AND password = $2`,
		email, password)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
