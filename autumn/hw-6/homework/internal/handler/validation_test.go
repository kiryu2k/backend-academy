package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework/internal/handler"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDevice(t *testing.T) {
	type Test struct {
		name string
		in   *handler.Device
		err  error
	}
	testCases := []Test{
		{
			name: "OK",
			in: &handler.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "1.1.1.1",
			},
		},
		{
			name: "empty serial num field",
			in: &handler.Device{
				Model: "some model",
				IP:    "1.1.1.1",
			},
			err: handler.ErrNotEnoughDeviceInfo,
		},
		{
			name: "empty model field",
			in: &handler.Device{
				SerialNum: "123",
				IP:        "1.1.1.1",
			},
			err: handler.ErrNotEnoughDeviceInfo,
		},
		{
			name: "empty ip field",
			in: &handler.Device{
				SerialNum: "123",
				Model:     "some model",
			},
			err: handler.ErrNotEnoughDeviceInfo,
		},
		{
			name: "too few octets",
			in: &handler.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "1.1.0",
			},
			err: handler.ErrInvalidIP,
		},
		{
			name: "too many octets",
			in: &handler.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "1.1.0.0.1",
			},
			err: handler.ErrInvalidIP,
		},
		{
			name: "no octets only points",
			in: &handler.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "...",
			},
			err: handler.ErrInvalidIP,
		},
		{
			name: "negative numbers",
			in: &handler.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "-15.1.1.1",
			},
			err: handler.ErrInvalidIP,
		},
		{
			name: "too large value for octet",
			in: &handler.Device{
				SerialNum: "123",
				Model:     "some model",
				IP:        "256.1.1.1",
			},
			err: handler.ErrInvalidIP,
		},
	}
	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := handler.ValidateDevice(test.in)
			assert.ErrorIs(t, err, test.err, test.name)
		})
	}
}

func BenchmarkValidateDevice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		octet := strconv.Itoa(i % 256)
		device := &handler.Device{
			SerialNum: strconv.Itoa(i),
			Model:     "some model",
			IP:        fmt.Sprintf("%s.%s.%s.%s", octet, octet, octet, octet),
		}
		_ = handler.ValidateDevice(device)
	}
}

func FuzzValidateDevice(f *testing.F) {
	testCases := []*handler.Device{
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
		device := new(handler.Device)
		if err := json.Unmarshal(data, device); err != nil {
			t.Skipf("invalid json data %s", string(data))
		}
		err := handler.ValidateDevice(device)
		expParsedIP := net.ParseIP(device.IP)
		if errors.Is(err, handler.ErrInvalidIP) && expParsedIP != nil {
			t.Errorf("got error %v validating device %#v", err, device)
		}
	})
}
