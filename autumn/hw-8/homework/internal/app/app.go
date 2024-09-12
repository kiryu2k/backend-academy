package app

import (
	"homework/internal/device"
)

type App struct {
	repo device.Repository
}

func New(repo device.Repository) *App {
	return &App{
		repo: repo,
	}
}

func (a *App) GetDevice(serialNum string) (*device.Device, error) {
	return a.repo.GetDevice(serialNum)
}

func (a *App) CreateDevice(d *device.Device) error {
	return a.repo.CreateDevice(d)
}

func (a *App) DeleteDevice(serialNum string) error {
	return a.repo.DeleteDevice(serialNum)
}

func (a *App) UpdateDevice(d *device.Device) error {
	return a.repo.UpdateDevice(d)
}
