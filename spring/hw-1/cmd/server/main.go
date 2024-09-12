package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"ws-chat/internal/adapters"
	"ws-chat/internal/adapters/pgrepo"
	"ws-chat/internal/transport/ws"
	"ws-chat/internal/usecase"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	if err := godotenv.Load(".env"); err != nil {
		logger.Warn(err.Error())
	}
	ctx := context.Background()
	connPool, err := pgrepo.NewConnectionPool(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func() {
		logger.Info("closing db connections...")
		connPool.Close()
	}()
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	errGroup := new(errgroup.Group)
	errGroup.Go(func() error {
		select {
		case s := <-sigChan:
			return errors.Errorf("captured signal: %v", s)
		}
	})
	var (
		msgRepo = adapters.NewMessageRepo(connPool)
		hub     = usecase.New(msgRepo, logger)
		server  = ws.New(os.Getenv("SERVER_ADDR"), hub, logger)
	)
	go func() {
		logger.Info("http server is starting...", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil {
			logger.Info(err.Error())
		}
	}()
	if err := errGroup.Wait(); err != nil {
		logger.Info("gracefully shutting down the server: " + err.Error())
	}
	if err := server.Shutdown(ctx); err != nil {
		logger.Info("failed to shutdown http server: " + err.Error())
	}
}
