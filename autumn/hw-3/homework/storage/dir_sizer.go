package storage

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{maxWorkersCount: 10}
}

func (s *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	dirs, files, err := d.Ls(ctx)
	if err != nil {
		return Result{}, err
	}
	res, err := calcFilesInfo(ctx, files)
	if err != nil {
		return Result{}, err
	}
	var (
		sigChan = make(chan struct{}, s.maxWorkersCount)
		resChan = make(chan Result, len(dirs))
	)
	errGroup, ctx := errgroup.WithContext(ctx)
	for _, dir := range dirs {
		sigChan <- struct{}{}
		dir := dir
		errGroup.Go(func() error {
			defer func() {
				<-sigChan
			}()
			r, err := s.Size(ctx, dir)
			if err != nil {
				return err
			}
			resChan <- r
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		return Result{}, err
	}
	close(resChan)
	for r := range resChan {
		res.Count += r.Count
		res.Size += r.Size
	}
	return res, nil
}

func calcFilesInfo(ctx context.Context, files []File) (Result, error) {
	res := Result{Count: int64(len(files))}
	for _, file := range files {
		size, err := file.Stat(ctx)
		if err != nil {
			return Result{}, err
		}
		res.Size += size
	}
	return res, nil
}
