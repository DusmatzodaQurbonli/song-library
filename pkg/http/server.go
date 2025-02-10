package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/DusmatzodaQurbonli/song-library/docs"
	"github.com/DusmatzodaQurbonli/song-library/internal/config"
	"github.com/DusmatzodaQurbonli/song-library/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
	log    *logrus.Logger
	config *config.Config
}

func NewServer(handler *handler.SongHandler, log *logrus.Logger, config *config.Config) *Server {
	router := gin.Default()

	server := &Server{
		router: router,
		log:    log,
		config: config,
	}

	router.Use(gin.Recovery())
	router.Use(server.loggingMiddleware)

	server.setupRoutes(handler)

	return server
}

func (s *Server) loggingMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()

	latency := time.Since(start)
	status := c.Writer.Status()

	s.log.Infof("[%d] %s %s | %s | %v",
		status, c.Request.Method, c.Request.URL.Path, c.ClientIP(), latency)
}

func (s *Server) setupRoutes(handler *handler.SongHandler) {
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := s.router.Group("/songs")
	{
		api.POST("/", handler.AddSong)
		api.GET("/:id", handler.GetSongText)
		api.GET("/", handler.GetSongs)
		api.PUT("/:id", handler.UpdateSong)
		api.DELETE("/:id", handler.DeleteSong)
	}
}

func (s *Server) Run() error {
	if s.config.Server.Host == "" || s.config.Server.Port == "" {
		s.log.Fatal("Server host or port is not configured")
	}

	server := &http.Server{
		Addr:    net.JoinHostPort(s.config.Server.Host, s.config.Server.Port),
		Handler: s.router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Errorf("Server error: %v", err)
		}
	}()

	s.log.Infof("Server is running on host %s and port %s", s.config.Server.Host, s.config.Server.Port)

	<-done
	s.log.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		s.log.Errorf("Server shutdown error: %v", err)
	}

	s.log.Info("Server stopped")
	return nil
}
