package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	port            int
	shutdownTimeout time.Duration
	DB              *pgxpool.Pool
}

func New(port int, shutdownTimeout time.Duration, db *pgxpool.Pool) *Server {
	return &Server{
		port:            port,
		shutdownTimeout: shutdownTimeout,
		DB:              db,
	}
}

func (s *Server) Serve() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		Handler:           s.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		BaseContext:       func(l net.Listener) context.Context { return ctx },
		ErrorLog:          log.New(os.Stderr, "http: ", log.LstdFlags),
	}

	errorCh := make(chan error, 1)
	go func() {
		log.Printf("Server running on port: %d", s.port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errorCh <- err
			return
		}
		errorCh <- nil
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}

		if s.DB != nil {
			s.DB.Close()
		}

		return nil
	case err := <-errorCh:
		if err != nil {
			return fmt.Errorf("listen error: %w", err)
		}
		return nil
	}
}
