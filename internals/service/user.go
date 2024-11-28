package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"

	"train_http_storage/internals/models"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserStorage interface {
	Create(ctx context.Context, user *models.User) error
	GetByCredentials(ctx context.Context, email, password string) (*models.User, error)
}

type Users struct {
	repo   UserStorage
	hasher PasswordHasher

	hmacSecret []byte
	tokenTll   time.Duration
}

func NewUserService(repo UserStorage, hasher PasswordHasher, secret []byte, tokenTll time.Duration) *Users {
	return &Users{
		repo:       repo,
		hasher:     hasher,
		hmacSecret: secret,
		tokenTll:   tokenTll,
	}
}

func (u *Users) SingUp(ctx context.Context, inp *models.SingUpInput) error {
	password, err := u.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return u.repo.Create(ctx, user)
}

func (u *Users) SingIn(ctx context.Context, userInput *models.SingInInput) (string, error) {
	password, err := u.hasher.Hash(userInput.Password)
	if err != nil {
		return "", err
	}

	user, err := u.repo.GetByCredentials(ctx, userInput.Email, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(u.tokenTll).Unix(),
	})

	return token.SignedString(u.hmacSecret)
}

func (u *Users) ParseToken(ctx context.Context, token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return u.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}
