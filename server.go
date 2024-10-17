package service_init_tool

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lie-flat-planet/service-init-tool/enum"
	"github.com/lie-flat-planet/service-init-tool/log"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Name string
	// service code
	Code int
	// LogLevel 默认是 debug 等级
	LogLevel string `env:""`
	HttpPort uint   `env:""`
	RunMode  string `env:""`

	httpServer *http.Server `skip:""`
}

func (s *Server) GetHttpPort() uint {
	return s.HttpPort
}

func (s *Server) Init() error {
	if err := s.check(); err != nil {
		return err
	}

	s.init()

	return nil
}

func (s *Server) init() {
	if s.LogLevel == "" {
		s.LogLevel = enum.Debug
	}

	if s.HttpPort == 0 {
		s.HttpPort = 80
	}

	if s.RunMode == "" {
		s.RunMode = enum.Debug
	}

	// log
	(&log.Log{
		Name:  s.Name,
		Level: s.LogLevel,
	}).SetDefaults().Build()
}

func (s *Server) check() error {
	if s.Name == "" {
		return fmt.Errorf("server Name cannot be empty")
	}

	if s.Code == 0 {
		return fmt.Errorf("server Code cannot be zero")
	}

	return nil
}

func (s *Server) GinServe(engine *gin.Engine) {
	port := fmt.Sprintf(":%d", s.HttpPort)
	s.httpServer = &http.Server{
		Addr:    port,
		Handler: engine,
	}

	//start server
	go func() {
		logrus.Println("Starting server on " + port + "...")
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("server startup error. err:%v", err)
		}
	}()

	// gracefulShutdown
	s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	logrus.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Println("Server exiting")
}
