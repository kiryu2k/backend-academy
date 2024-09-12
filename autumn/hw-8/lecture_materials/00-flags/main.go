package main

import (
	"flag"
	"fmt"
)

// go run main.go --port=8090

func main() {
	host := flag.String("host", "", "server host")
	port := flag.Int("port", 0, "server port")
	logLevel := flag.String("log-level", "", "server log level (INFO, ERROR, DEBUG)")

	flag.Parse()

	fmt.Printf("host: %s, port: %d, log level: %s\n", *host, *port, *logLevel)
}
