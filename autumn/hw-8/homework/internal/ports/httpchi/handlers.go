package httpchi

import (
	"encoding/json"
	"homework/internal/device"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func createDevice(usecase device.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		d := new(device.Device)
		if err := json.NewDecoder(r.Body).Decode(d); err != nil {
			writeErrorResponse(w, ErrInvalidDataFormat, http.StatusBadRequest)
			return
		}
		if err := Validate(d); err != nil {
			writeErrorResponse(w, err, http.StatusBadRequest)
			return
		}
		if err := usecase.CreateDevice(d); err != nil {
			writeErrorResponse(w, err, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getDevice(usecase device.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		serialNum := chi.URLParam(r, "serialNum")
		device, err := usecase.GetDevice(serialNum)
		if err != nil {
			writeErrorResponse(w, err, http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(device); err != nil {
			writeErrorResponse(w, err, http.StatusInternalServerError)
		}
	}
}

func updateDevice(usecase device.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		d := new(device.Device)
		if err := json.NewDecoder(r.Body).Decode(d); err != nil {
			writeErrorResponse(w, ErrInvalidDataFormat, http.StatusBadRequest)
			return
		}
		if err := Validate(d); err != nil {
			writeErrorResponse(w, err, http.StatusBadRequest)
			return
		}
		if err := usecase.UpdateDevice(d); err != nil {
			writeErrorResponse(w, err, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func deleteDevice(usecase device.Usecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		serialNum := chi.URLParam(r, "serialNum")
		if err := usecase.DeleteDevice(serialNum); err != nil {
			writeErrorResponse(w, err, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
