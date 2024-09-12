package repository

import (
	"errors"
)

var (
	ErrDeviceNotFound = errors.New("specified device is not found")
	ErrDeviceExists   = errors.New("specified device already exists")
)
