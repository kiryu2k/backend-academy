package main

import (
	"flag"
	"homework/internal/config"
	"homework/internal/repository"
	"homework/internal/transport/grpcserver"
	"homework/internal/usecase"
	"log"
	"time"
)

func main() {
	cfgPath := flag.String("cfg", "./configs/config.yaml", "path to yaml config file")
	flag.Parse()
	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	fileRepo, err := repository.New(cfg.Dirpath)
	if err != nil {
		log.Fatal(err)
	}
	usecase := usecase.New(fileRepo)
	server := grpcserver.New(cfg.Addr, usecase)
	go func() {
		for {
			if err := fileRepo.Update(); err != nil {
				log.Printf("error updating file repository: %v", err)
			}
			/* каждые 30 секунд обновляем мапу с данными о файлах */
			time.Sleep(30 * time.Second)
		}
	}()
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
