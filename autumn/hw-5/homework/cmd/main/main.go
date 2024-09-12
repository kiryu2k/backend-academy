package main

import (
	"homework/internal/handler"
	"homework/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type serverCfg struct {
	addr string
}

func main() {
	var (
		s   = service.New()
		h   = handler.New(s)
		mux = chi.NewMux()
		cfg = NewServerCfg()
	)
	mux.Post("/device", h.CreateDevice)
	mux.Get("/device/{serialNum}", h.GetDevice)
	mux.Put("/device", h.UpdateDevice)
	mux.Delete("/device/{serialNum}", h.DeleteDevice)
	server := http.Server{
		Addr:    cfg.addr,
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func NewServerCfg() *serverCfg {
	addr := os.Getenv("HTTP_SERVER_ADDR")
	port := os.Getenv("HTTP_SERVER_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return &serverCfg{
		addr: addr + ":" + port,
	}
}
