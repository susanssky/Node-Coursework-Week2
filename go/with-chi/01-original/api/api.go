package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"node-week-02-with-chi/handlers"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type APIServer struct {
	Addr    string
	Handler *handlers.MessageHandler
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		Addr:    addr,
		Handler: handlers.New(),
	}
}

func (s *APIServer) Run() error {
	router := s.Routes()

	fmt.Printf("Server starting on port %v\n", s.Addr)

	srv := http.Server{
		Addr: s.Addr, Handler: router,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Error:%v\n", err)
		}
	}()

	fmt.Println("Press Ctrl+C to stop the server")
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("error: %v\n", err)
	}

	fmt.Println("Server gracefully stopped")

	return nil
}
