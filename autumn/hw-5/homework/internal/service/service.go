package service

import (
	"fmt"
	"homework/internal/handler"
	"sync"
)

type Service struct {
	devices map[string]handler.Device
	mu      sync.RWMutex
}

func New() *Service {
	return &Service{
		devices: make(map[string]handler.Device),
	}
}

func (s *Service) GetDevice(serialNum string) (handler.Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	device, ok := s.devices[serialNum]
	if !ok {
		return handler.Device{}, fmt.Errorf("%w: serial number '%s'", errDeviceNotFound, serialNum)
	}
	return device, nil
}

func (s *Service) CreateDevice(d handler.Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.devices[d.SerialNum]; ok {
		return fmt.Errorf("%w: serial number '%s'", errDeviceExists, d.SerialNum)
	}
	s.devices[d.SerialNum] = d
	return nil
}

func (s *Service) DeleteDevice(serialNum string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.devices[serialNum]; !ok {
		return fmt.Errorf("%w: serial number '%s'", errDeviceNotFound, serialNum)
	}
	delete(s.devices, serialNum)
	return nil
}

func (s *Service) UpdateDevice(d handler.Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.devices[d.SerialNum]; !ok {
		return fmt.Errorf("%w: serial number '%s'", errDeviceNotFound, d.SerialNum)
	}
	s.devices[d.SerialNum] = d
	return nil
}
