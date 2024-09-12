package usecase

import (
	"context"
	"fmt"
	"homework/internal/domain"
)

type File struct {
	repo domain.FileRepo
}

func New(repo domain.FileRepo) *File {
	return &File{
		repo: repo,
	}
}

func (f *File) Get(filename string) ([]byte, error) {
	file, err := f.repo.Find(context.Background(), filename)
	if err != nil {
		return nil, fmt.Errorf("repo find file %s: %w", filename, err)
	}
	return file.Data, nil
}

func (f *File) All(ctx context.Context) []string {
	return f.repo.All(ctx)
}

func (f *File) GetInfo(ctx context.Context, filename string) (*domain.FileInfo, error) {
	file, err := f.repo.Find(ctx, filename)
	if err != nil {
		return nil, fmt.Errorf("repo find file %s: %w", filename, err)
	}
	return file, nil
}
