package repository

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"os"
	"path/filepath"
	"sync"
)

type fileRepo struct {
	mu      sync.RWMutex
	files   map[string]*domain.FileInfo
	dirname string
}

func New(dirname string) (*fileRepo, error) {
	files, err := loadFilesIntoMap(dirname)
	if err != nil {
		return nil, fmt.Errorf("load files into map: %w", err)
	}
	return &fileRepo{
		files:   files,
		dirname: dirname,
	}, nil
}

func (f *fileRepo) Find(ctx context.Context, filename string) (*domain.FileInfo, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	file, ok := f.files[filename]
	if !ok {
		return nil, domain.ErrFileNotFound
	}
	return file, nil
}

func (f *fileRepo) All(ctx context.Context) []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	result := make([]string, 0, len(f.files))
	for filename := range f.files {
		result = append(result, filename)
	}
	return result
}

func (f *fileRepo) Update() error {
	files, err := loadFilesIntoMap(f.dirname)
	if err != nil {
		return fmt.Errorf("load files into map: %w", err)
	}
	f.mu.Lock()
	f.files = files
	f.mu.Unlock()
	return nil
}

func loadFilesIntoMap(dirname string) (map[string]*domain.FileInfo, error) {
	/* open directory */
	dir, err := os.Open(dirname)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", dirname, err)
	}
	defer dir.Close()
	/* get list of all files in dir */
	list, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}
	files := make(map[string]*domain.FileInfo)
	for _, f := range list {
		if f.IsDir() {
			continue
		}
		data, err := os.ReadFile(dirname + f.Name())
		if err != nil {
			return nil, fmt.Errorf("read file %s: %w", f.Name(), err)
		}
		files[f.Name()] = &domain.FileInfo{
			Name: f.Name(),
			Data: data,
			Size: uint64(f.Size() / 1024), // размер в Кб
			Type: filepath.Ext(f.Name()),
		}
	}
	return files, nil
}
