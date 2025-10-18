package main

import (
	"flag"
	"log"
	"time"

	"github.com/cateruu/moto-backend/internals/server"
)

func main() {
	var port int
	var shutdownTimeout time.Duration

	flag.IntVar(&port, "port", 1337, "API server port")
	flag.DurationVar(&shutdownTimeout, "shutdown-timeout", 10*time.Second, "API server shutdown timeout")

	flag.Parse()

	srv := server.New(
		port,
		shutdownTimeout,
	)

	if err := srv.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
