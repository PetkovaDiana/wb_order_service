package server

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	Port         string `yaml:"port"`
	ReadTimeout  int64  `yaml:"read_timeout"`
	WriteTimeout int64  `yaml:"write_timeout"`
}

type Server struct {
	httpServer *http.Server
	cfg        *Config
}

func NewServer(cfg *Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + s.cfg.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    s.getReadTimeout(),
		WriteTimeout:   s.getWriteTimeout(),
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) getReadTimeout() time.Duration {
	return time.Duration(s.cfg.ReadTimeout) * time.Millisecond
}

func (s *Server) getWriteTimeout() time.Duration {
	return time.Duration(s.cfg.WriteTimeout) * time.Millisecond
}
