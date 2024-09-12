package main

import (
	"flag"
	"homework/internal/adapters/devicerepo"
	"homework/internal/app"
	"homework/internal/ports/httpchi"
	"log"

	srvrcfg "github.com/kiryu-dev/server-config"
)

func main() {
	cfgPath := flag.String("config", "", "http server configuration file")
	flag.Parse()
	if cfgPath == nil || *cfgPath == "" {
		log.Fatal("config file is not specified")
	}
	cfg, err := srvrcfg.LoadYamlCfg(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	var (
		repo    = devicerepo.New()
		usecase = app.New(repo)
		server  = httpchi.NewHTTPServer(cfg, usecase)
	)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
