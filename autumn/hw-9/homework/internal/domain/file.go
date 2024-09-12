package domain

import (
	"context"
	"errors"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

type FileInfo struct {
	Name string
	Data []byte
	Size uint64
	Type string
}

//go:generate mockery --name FileUseCase
type FileUseCase interface {
	Get(string) ([]byte, error)
	All(context.Context) []string
	GetInfo(context.Context, string) (*FileInfo, error)
}

//go:generate mockery --name FileRepo
type FileRepo interface {
	Find(context.Context, string) (*FileInfo, error)
	All(context.Context) []string
}
