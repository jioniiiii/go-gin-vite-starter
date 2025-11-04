package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type AppConfig struct {
	Env            string
	GinMode        string
	Port           string
	AllowOrigins   []string
	TrustedProxies []string

	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration

	DatabaseURL string
	JWTSecret   string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getdur(key, def string) time.Duration {
	v := getenv(key, def)
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("invalid duration for %s=%q, using default %s", key, v, def)
		d, _ = time.ParseDuration(def)
	}
	return d
}

func splitCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func LoadConfig() AppConfig {
	cfg := AppConfig{
		Env:            getenv("ENV", "production"),
		GinMode:        getenv("GIN_MODE", "release"),
		Port:           getenv("PORT", "8080"),
		AllowOrigins:   splitCSV(getenv("ALLOW_ORIGINS", "")),
		TrustedProxies: splitCSV(getenv("TRUSTED_PROXIES", "0.0.0.0/0")),

		ReadTimeout:  getdur("READ_TIMEOUT", "10s"),
		WriteTimeout: getdur("WRITE_TIMEOUT", "15s"),
		IdleTimeout:  getdur("IDLE_TIMEOUT", "60s"),

		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	// Default to common origin
	if cfg.Env == "development" {
		if len(cfg.AllowOrigins) == 0 {
			cfg.AllowOrigins = []string{"http://localhost:5173"}
		}
	}

	// Basic production safety checks
	if cfg.Env == "production" {
		if len(cfg.AllowOrigins) == 0 {
			log.Println("WARNING: ALLOW_ORIGINS is empty in production. Set it to your frontend domain(s).")
		}
		if cfg.JWTSecret == "" {
			log.Println("WARNING: JWT_SECRET is empty in production. Set a strong secret.")
		}
	}

	return cfg
}

func ApplyGinMode(cfg AppConfig) {
	gin.SetMode(cfg.GinMode) // release/debug/test
}

func ApplyTrustedProxies(r *gin.Engine, cfg AppConfig) {
	// cfg.TrustedProxies is a []string of CIDRs/IPs (e.g. "10.0.0.0/8", "127.0.0.1")
	if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		log.Printf("SetTrustedProxies error: %v", err)
	}
}

func Cors(cfg AppConfig) gin.HandlerFunc {
	allow := cfg.AllowOrigins
	// For prod
	return cors.New(cors.Config{
		AllowOrigins:     allow,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
