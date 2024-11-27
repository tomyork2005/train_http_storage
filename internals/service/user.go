package service

import (
	"context"
	"time"
	"train_http_storage/internals/models"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserStorage interface {
	Create(ctx context.Context, user models.User) error
}

type Users struct {
	repo   UserStorage
	hasher PasswordHasher
}

func NewUsers(repo UserStorage, hasher PasswordHasher) *Users {
	return &Users{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *Users) SingUp(ctx context.Context, inp models.SingUpInput) error {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return s.repo.Create(ctx, user)
}
