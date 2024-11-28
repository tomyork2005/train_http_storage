package service

import (
	"context"
	"train_http_storage/internals/models"
)

type FileStorage interface {
	GetAll(ctx context.Context, id int64) ([]models.File, error)
}

type FileService struct {
	fileStorage FileStorage
}

func NewFileService(storage FileStorage) *FileService {
	return &FileService{
		fileStorage: storage,
	}
}

func (f *FileService) GetAll(ctx context.Context, id int64) ([]models.File, error) {
	return f.fileStorage.GetAll(ctx, id)
}
