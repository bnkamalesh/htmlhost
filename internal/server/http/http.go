package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bnkamalesh/htmlhost/internal/api"
	"github.com/bnkamalesh/webgo/v5"
	"github.com/bnkamalesh/webgo/v5/middleware/accesslog"
)

type Server struct {
	cfg    *Config
	router *webgo.Router
}

func (h *Server) Start() error {
	if h.router == nil {
		return errors.New("router not initialized")
	}

	host := h.cfg.Host

	if len(h.cfg.Port) > 0 {
		host = fmt.Sprintf("%s:%s", host, h.cfg.Port)
	}

	httpServer := &http.Server{
		Addr:           host,
		Handler:        h.router,
		ReadTimeout:    h.cfg.ReadTimeout,
		WriteTimeout:   h.cfg.WriteTimeout,
		MaxHeaderBytes: h.cfg.MaxHeaderSize,
	}

	webgo.LOGHANDLER.Info("HTTP server, listening on", host)

	return httpServer.ListenAndServe()
}

type Config struct {
	Host             string
	Port             string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	MaxHeaderSize    int
	MaxBodysizeBytes int

	GeneratedBaseURL string
}

func (c *Config) Sanitize() {
	c.Host = strings.TrimSpace(c.Host)
	c.Port = strings.TrimSpace(c.Port)
	if c.Port == "" {
		c.Port = "8000"
	}
	if c.ReadTimeout < time.Second {
		c.ReadTimeout = time.Second * 5
	}
	if c.WriteTimeout < time.Second {
		c.WriteTimeout = time.Second * 5
	}

	if c.MaxHeaderSize < 256 {
		// 10KB
		c.MaxHeaderSize = 1024 * 10
	}

	if c.MaxBodysizeBytes < 256 {
		// 1MB
		c.MaxBodysizeBytes = 1024 * 1024 * 2
	}
}

func New(cfg *Config, api *api.API) (*Server, error) {
	cfg.Sanitize()

	wcfg := &webgo.Config{
		Host:         cfg.Host,
		Port:         cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	handlers, err := newHandler(api, cfg.GeneratedBaseURL)
	if err != nil {
		return nil, err
	}

	httproutes := routes(handlers)
	router := webgo.NewRouter(wcfg, httproutes)
	router.Use(accesslog.AccessLog)
	router.UseOnSpecialHandlers(accesslog.AccessLog)

	return &Server{
		cfg:    cfg,
		router: router,
	}, nil
}
