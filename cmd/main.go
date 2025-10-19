package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/cateruu/moto-backend/internals/db"
	"github.com/cateruu/moto-backend/internals/server"
)

func main() {
	var port int
	var shutdownTimeout time.Duration

	flag.IntVar(&port, "port", 1337, "API server port")
	flag.DurationVar(&shutdownTimeout, "shutdown-timeout", 10*time.Second, "API server shutdown timeout")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := db.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	srv := server.New(
		port,
		shutdownTimeout,
		db,
	)

	if err := srv.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
