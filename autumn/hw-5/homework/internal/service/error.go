package service

import (
	"errors"
)

var (
	errDeviceNotFound = errors.New("specified device is not found")
	errDeviceExists   = errors.New("specified device already exists")
)
