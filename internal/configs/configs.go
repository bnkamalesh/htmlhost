package configs

import (
	"os"
	"strings"
	"time"

	"github.com/bnkamalesh/htmlhost/internal/pages"
	"github.com/bnkamalesh/htmlhost/internal/server/http"
)

type Configs struct{}

func (cfg *Configs) getEnv(key string) string {
	return strings.Trim(strings.TrimSpace(os.Getenv(key)), "\n")
}

func (cfg *Configs) HTTP() *http.Config {
	baseURL := cfg.getEnv("GENERATED_BASEURL")
	if baseURL == "" {
		baseURL = "https://htmlhost.live"
	}

	return &http.Config{
		Host:             cfg.getEnv("HTTP_HOST"),
		Port:             cfg.getEnv("HTTP_PORT"),
		ReadTimeout:      time.Second * 3,
		WriteTimeout:     time.Hour * 1,
		GeneratedBaseURL: baseURL,
	}
}

func (cfg *Configs) Pages() *pages.Config {
	host := cfg.getEnv("DATASTORE_HOST")
	if host == "" {
		host = "localhost"
	}

	port := cfg.getEnv("DATASTORE_PORT")
	if port == "" {
		port = "6379"
	}

	return &pages.Config{
		Host:             host,
		Port:             port,
		StoreName:        cfg.getEnv("DATASTORE_NAME"),
		Password:         cfg.getEnv("DATASTORE_PASSWORD"),
		PoolSize:         25,
		IdleTimeoutSecs:  time.Second * 60,
		ReadTimeoutSecs:  time.Second * 5,
		WriteTimeoutSecs: time.Second * 5,
		DialTimeoutSecs:  time.Second * 1,
	}
}

func New() *Configs {
	return &Configs{}
}
