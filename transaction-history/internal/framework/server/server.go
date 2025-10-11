package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"transaction-history/internal/framework/logger"

	"github.com/labstack/echo/v4"
)

type Server struct {
	name        string
	gracePeriod time.Duration

	e *echo.Echo
	h *http.Server
}

type Config struct {
	Name         string
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	GracePeriod  time.Duration
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func Serve(e *echo.Echo, c *Config) *Server {
	log := logger.GetLogger(c.Name)
	defer log.Sync()

	log.Infow("start server", "name", c.Name, "host", c.Host, "port", c.Port)
	h := &http.Server{
		Addr:         c.Addr(),
		Handler:      e,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	}
	go func() {
		if err := h.ListenAndServe(); err != http.ErrServerClosed {
			log.Panic(err)
		}
	}()

	return &Server{
		name:        c.Name,
		gracePeriod: c.GracePeriod,
		e:           e,
		h:           h,
	}
}

func (s *Server) Shutdown() {
	log := logger.GetLogger(s.name)
	defer log.Sync()

	log.Infow("shutdown server", "name", s.name)

	ctx, cancelContext := context.WithTimeout(context.TODO(), s.gracePeriod)
	defer cancelContext()
	if err := s.h.Shutdown(ctx); err != nil {
		log.Error(fmt.Sprintf("Could not gracefully shut down web server. Error: %s", err.Error()))
	}
}
