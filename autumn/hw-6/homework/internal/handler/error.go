package handler

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrInvalidDataFormat   = errors.New("invalid input data format")
	ErrInvalidIP           = errors.New("invalid ipv4 address format")
	ErrNotEnoughDeviceInfo = errors.New("all device info fields are required")
)

type errResponse struct {
	Msg string `json:"message"`
}

func writeErrorResponse(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(&errResponse{Msg: err.Error()})
}
