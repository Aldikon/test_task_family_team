package service

import (
	"context"
	"fmt"
	"os"
)

type FileRepo interface {
	Get(ctx context.Context) []string
	GetByName(ctx context.Context, name string) (*os.File, error)
}

type file struct {
	fileRepo FileRepo
}

func NewFile(fileRepo FileRepo) *file {
	return &file{fileRepo: fileRepo}
}

func (s *file) Get(ctx context.Context) []string {
	return s.fileRepo.Get(ctx)
}

func (s *file) GetByName(ctx context.Context, name string) (*os.File, error) {
	f, err := s.fileRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("service 'file', func GetByName, err: %w", err)
	}

	return f, nil
}
