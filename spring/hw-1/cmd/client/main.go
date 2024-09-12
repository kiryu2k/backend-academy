package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
	u := url.URL{Scheme: "ws", Host: os.Getenv("SERVER_ADDR"), Path: "/"}
	username := ""
	for username == "" {
		fmt.Print("enter your name: ")
		_, _ = fmt.Scan(&username)
		username = strings.TrimSpace(username)
	}
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	errGroup := new(errgroup.Group)
	errGroup.Go(func() error {
		select {
		case s := <-sigChan:
			return errors.Errorf("captured signal: %v", s)
		}
	})
	go func() {
		logger.Info("connecting to " + u.String())
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), map[string][]string{
			"X-User-Name-Key": {username},
		})
		if err != nil {
			logger.Fatal("dial: " + err.Error())
		}
		defer func() {
			_ = conn.Close()
		}()
		go readMessages(conn, logger)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := scanner.Bytes()
			fmt.Printf("\033[1A\033[K")
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logger.Fatal("write: " + err.Error())
			}
		}
	}()
	if err := errGroup.Wait(); err != nil {
		logger.Info("gracefully stopping: " + err.Error())
	}
}

func readMessages(conn *websocket.Conn, logger *zap.Logger) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Warn("read: " + err.Error())
			return
		}
		fmt.Printf("%s\n", message)
	}
}
