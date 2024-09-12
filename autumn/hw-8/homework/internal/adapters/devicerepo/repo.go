package devicerepo

import (
	"fmt"
	"homework/internal/device"
	"sync"
)

type repo struct {
	devices map[string]*device.Device
	mu      sync.RWMutex
}

func New() *repo {
	return &repo{
		devices: make(map[string]*device.Device),
	}
}

func (r *repo) GetDevice(serialNum string) (*device.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.devices[serialNum]
	if !ok {
		return &device.Device{}, fmt.Errorf("%w: serial number '%s'", ErrDeviceNotFound, serialNum)
	}
	return d, nil
}

func (r *repo) CreateDevice(d *device.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.devices[d.SerialNum]; ok {
		return fmt.Errorf("%w: serial number '%s'", ErrDeviceExists, d.SerialNum)
	}
	r.devices[d.SerialNum] = d
	return nil
}

func (r *repo) DeleteDevice(serialNum string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.devices[serialNum]; !ok {
		return fmt.Errorf("%w: serial number '%s'", ErrDeviceNotFound, serialNum)
	}
	delete(r.devices, serialNum)
	return nil
}

func (r *repo) UpdateDevice(d *device.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.devices[d.SerialNum]; !ok {
		return fmt.Errorf("%w: serial number '%s'", ErrDeviceNotFound, d.SerialNum)
	}
	r.devices[d.SerialNum] = d
	return nil
}
