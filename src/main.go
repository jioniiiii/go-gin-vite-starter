package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1) Load config from env
	cfg := LoadConfig()
	ApplyGinMode(cfg)

	// 2) Gin setup
	r := gin.New()
	r.Use(gin.Recovery(), Cors(cfg))
	ApplyTrustedProxies(r, cfg)

	// 3) Health + api
	r.GET("/api/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// 4) Timeouts + graceful shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Printf("listening on :%s (env=%s, mode=%s)", cfg.Port, cfg.Env, cfg.GinMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
