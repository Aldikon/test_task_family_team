package repository

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/Aldikon/test_task_family_team/internal/model"
)

type file struct {
	files    []string
	fs       sync.Map
	destPath string
}

func NewFile(fsDir []fs.DirEntry, destPath string) *file {
	fileRepo := file{
		files:    make([]string, 0, len(fsDir)),
		destPath: destPath,
	}

	for _, f := range fsDir {
		fileRepo.files = append(fileRepo.files, f.Name())
		fileRepo.fs.Store(f.Name(), f)
	}

	return &fileRepo
}

func (r *file) Get(ctx context.Context) []string {
	dst := make([]string, len(r.files))
	copy(dst, r.files)
	return dst
}

func (r *file) GetByName(ctx context.Context, name string) (*os.File, error) {
	val, ok := r.fs.Load(name)
	if !ok {
		return nil, model.ErrNotFound
	}

	v, ok := val.(fs.DirEntry)
	if !ok {
		return nil, fmt.Errorf("cast file to fs.DirEntry, name: %s", name)
	}

	f, err := os.Open(filepath.Join(r.destPath, v.Name()))
	if err != nil {
		return nil, fmt.Errorf("open file, err: %w", err)
	}

	return f, nil
}
