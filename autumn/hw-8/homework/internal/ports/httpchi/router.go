package httpchi

import (
	"homework/internal/device"

	"github.com/go-chi/chi/v5"
)

func AppRouter(mux *chi.Mux, usecase device.Usecase) {
	mux.Post("/device", createDevice(usecase))
	mux.Get("/device/{serialNum}", getDevice(usecase))
	mux.Put("/device", updateDevice(usecase))
	mux.Delete("/device/{serialNum}", deleteDevice(usecase))
}
