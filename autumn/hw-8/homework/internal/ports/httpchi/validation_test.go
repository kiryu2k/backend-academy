package httpchi_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework/internal/device"
	"homework/internal/ports/httpchi"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDevice(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		in   *device.Device
		err  error
	}{
		{
			name: "OK",
			in: &device.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "1.1.1.1",
			},
		},
		{
			name: "empty serial num field",
			in: &device.Device{
				Model: "some model",
				IP:    "1.1.1.1",
			},
			err: httpchi.ErrNotEnoughDeviceInfo,
		},
		{
			name: "empty model field",
			in: &device.Device{
				SerialNum: "123",
				IP:        "1.1.1.1",
			},
			err: httpchi.ErrNotEnoughDeviceInfo,
		},
		{
			name: "empty ip field",
			in: &device.Device{
				SerialNum: "123",
				Model:     "some model",
			},
			err: httpchi.ErrNotEnoughDeviceInfo,
		},
		{
			name: "too few octets",
			in: &device.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "1.1.0",
			},
			err: httpchi.ErrInvalidIP,
		},
		{
			name: "too many octets",
			in: &device.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "1.1.0.0.1",
			},
			err: httpchi.ErrInvalidIP,
		},
		{
			name: "no octets only points",
			in: &device.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "...",
			},
			err: httpchi.ErrInvalidIP,
		},
		{
			name: "negative numbers",
			in: &device.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "-15.1.1.1",
			},
			err: httpchi.ErrInvalidIP,
		},
		{
			name: "too large value for octet",
			in: &device.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "256.1.1.1",
			},
			err: httpchi.ErrInvalidIP,
		},
	}
	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := httpchi.Validate(test.in)
			assert.ErrorIs(t, err, test.err, test.name)
		})
	}
}

func BenchmarkValidateDevice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		octet := strconv.Itoa(i % 256)
		device := &device.Device{
			SerialNum: strconv.Itoa(i),
			Model:     "some model",
			IP:        fmt.Sprintf("%s.%s.%s.%s", octet, octet, octet, octet),
		}
		_ = httpchi.Validate(device)
	}
}

func FuzzValidateDevice(f *testing.F) {
	testCases := []*device.Device{
		{
			SerialNum: "16",
			Model:     "some model",
			IP:        "1.1.1.1",
		},
		{
			SerialNum: "888",
			Model:     "model123213",
			IP:        "192.32.0.1",
		},
		{
			SerialNum: "333999",
			Model:     "brand new model",
			IP:        "255.255.8.16",
		},
	}
	for _, test := range testCases {
		data, _ := json.Marshal(test)
		f.Add(data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if !json.Valid(data) {
			t.Skipf("invalid json data %s", string(data))
		}
		device := new(device.Device)
		if err := json.Unmarshal(data, device); err != nil {
			t.Skipf("invalid json data %s", string(data))
		}
		err := httpchi.Validate(device)
		expParsedIP := net.ParseIP(device.IP)
		if errors.Is(err, httpchi.ErrInvalidIP) && expParsedIP != nil {
			t.Errorf("got error %v validating device %#v", err, device)
		}
	})
}
