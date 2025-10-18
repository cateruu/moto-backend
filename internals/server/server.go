package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	port            int
	shutdownTimeout time.Duration
}

func New(port int, shutdownTimeout time.Duration) *Server {
	return &Server{
		port:            port,
		shutdownTimeout: shutdownTimeout,
	}
}

func (s *Server) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
		return nil
	case err := <-errorCh:
		if err != nil {
			return fmt.Errorf("listen error: %w", err)
		}
		return nil
	}
}
