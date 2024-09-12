package handler

import (
	"fmt"
	"strconv"
	"strings"
)

func validateDevice(d *Device) error {
	if d.SerialNum == "" {
		return fmt.Errorf("%w: 'serial_num' cannot be empty", errNotEnoughDeviceInfo)
	}
	if d.Model == "" {
		return fmt.Errorf("%w: 'model' cannot be empty", errNotEnoughDeviceInfo)
	}
	return validateIp(d.IP)
}

func validateIp(ip string) error {
	if ip == "" {
		return fmt.Errorf("%w: 'ip' cannot be empty", errNotEnoughDeviceInfo)
	}
	octets := strings.Split(ip, ".")
	if len(octets) != 4 {
		return errInvalidIP
	}
	for _, octet := range octets {
		if _, err := strconv.ParseUint(octet, 10, 8); err != nil {
			return fmt.Errorf("%w: %w", errInvalidIP, err)
		}
	}
	return nil
}
