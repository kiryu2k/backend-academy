package main

import (
	"fmt"
	"os"
)

// export APP_PORT=8080
// printenv

func main() {
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	logLevel := os.Getenv("APP_LOG_LEVEL")

	fmt.Printf("host: %s, port: %s, log level: %s\n", host, port, logLevel)
}
