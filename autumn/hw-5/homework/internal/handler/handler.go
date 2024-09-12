package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

type Device struct {
	SerialNum string `json:"serial_num"`
	Model     string `json:"model"`
	IP        string `json:"ip"`
}

type Service interface {
	GetDevice(string) (Device, error)
	CreateDevice(Device) error
	DeleteDevice(string) error
	UpdateDevice(Device) error
}

func New(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d := new(Device)
	if err := json.NewDecoder(r.Body).Decode(d); err != nil {
		writeErrorResponse(w, errInvalidDataFormat, http.StatusBadRequest)
		return
	}
	if err := validateDevice(d); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if err := h.service.CreateDevice(*d); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	serialNum := chi.URLParam(r, "serialNum")
	device, err := h.service.GetDevice(serialNum)
	if err != nil {
		writeErrorResponse(w, err, http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(&device); err != nil {
		writeErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d := new(Device)
	if err := json.NewDecoder(r.Body).Decode(d); err != nil {
		writeErrorResponse(w, errInvalidDataFormat, http.StatusBadRequest)
		return
	}
	if err := validateDevice(d); err != nil {
		writeErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateDevice(*d); err != nil {
		writeErrorResponse(w, err, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	serialNum := chi.URLParam(r, "serialNum")
	if err := h.service.DeleteDevice(serialNum); err != nil {
		writeErrorResponse(w, err, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
