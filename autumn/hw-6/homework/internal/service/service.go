package service

import (
	"homework/internal/handler"
)

type Service struct {
	repo Repository
}

type Repository interface {
	GetDevice(string) (handler.Device, error)
	CreateDevice(handler.Device) error
	DeleteDevice(string) error
	UpdateDevice(handler.Device) error
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetDevice(serialNum string) (handler.Device, error) {
	return s.repo.GetDevice(serialNum)
}

func (s *Service) CreateDevice(d handler.Device) error {
	return s.repo.CreateDevice(d)
}

func (s *Service) DeleteDevice(serialNum string) error {
	return s.repo.DeleteDevice(serialNum)
}

func (s *Service) UpdateDevice(d handler.Device) error {
	return s.repo.UpdateDevice(d)
}
