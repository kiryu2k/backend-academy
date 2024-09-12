package httpchi

import (
	"homework/internal/device"
	"net/http"

	"github.com/go-chi/chi/v5"
	srvcfg "github.com/kiryu-dev/server-config"
)

func NewHTTPServer(cfg *srvcfg.ServerConfig, usecase device.Usecase) *http.Server {
	mux := chi.NewMux()
	AppRouter(mux, usecase)
	return &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
