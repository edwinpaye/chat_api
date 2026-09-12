package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat_api/internal/config"
	"chat_api/internal/platform"
	"chat_api/internal/service"
	"chat_api/internal/security"
	"chat_api/internal/transport/ws"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := platform.OpenDB(cfg.Database.URL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	tokens := security.NewTokenService(cfg.Security.JWTSecret, cfg.Security.JWTTTL())
	services := service.Build(db)
	hub := ws.NewHub()
	go hub.Run()

	handler := ws.NewHandler(hub, services, tokens, cfg)

	mux := http.NewServeMux()
	mux.HandleFunc(cfg.Server.WSPath, handler.ServeWS)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv := &http.Server{Addr: cfg.Server.Addr(), Handler: mux}

	go func() {
		scheme := "ws"
		var err error
{
			log.Printf("listening on %s://%s%s", scheme, cfg.Server.Addr(), cfg.Server.WSPath)
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("shutdown complete")
}