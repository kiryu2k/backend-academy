package repository

import (
	"fmt"
	"homework/internal/handler"
	"sync"
)

type Repository struct {
	devices map[string]handler.Device
	mu      sync.RWMutex
}

func New() *Repository {
	return &Repository{
		devices: make(map[string]handler.Device),
	}
}

func (r *Repository) GetDevice(serialNum string) (handler.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	device, ok := r.devices[serialNum]
	if !ok {
		return handler.Device{}, fmt.Errorf("%w: serial number '%s'", ErrDeviceNotFound, serialNum)
	}
	return device, nil
}

func (r *Repository) CreateDevice(d handler.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.devices[d.SerialNum]; ok {
		return fmt.Errorf("%w: serial number '%s'", ErrDeviceExists, d.SerialNum)
	}
	r.devices[d.SerialNum] = d
	return nil
}

func (r *Repository) DeleteDevice(serialNum string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.devices[serialNum]; !ok {
		return fmt.Errorf("%w: serial number '%s'", ErrDeviceNotFound, serialNum)
	}
	delete(r.devices, serialNum)
	return nil
}

func (r *Repository) UpdateDevice(d handler.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.devices[d.SerialNum]; !ok {
		return fmt.Errorf("%w: serial number '%s'", ErrDeviceNotFound, d.SerialNum)
	}
	r.devices[d.SerialNum] = d
	return nil
}
