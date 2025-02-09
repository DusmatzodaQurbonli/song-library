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

// Server struct
type Server struct {
	router *gin.Engine
	log    *logrus.Logger
	config *config.Config
}

func NewServer(handler *handler.SongHandler, log *logrus.Logger, config *config.Config) *Server {
	router := gin.Default()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		log.Infof("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	router.POST("/songs", handler.AddSong)
	router.GET("/songs/:id", handler.GetSongText)
	router.GET("/songs", handler.GetSongs)
	router.PUT("/songs/:id", handler.UpdateSong)
	router.DELETE("/songs/:id", handler.DeleteSong)

	return &Server{
		router: router,
		log:    log,
		config: config,
	}
}

// Run запускает сервер
func (s *Server) Run() error {
	if s.config.Server.Host == "" || s.config.Server.Port == "" {
		s.log.Fatal("Server host or port is not configured")
	}

	server := &http.Server{
		Addr:    net.JoinHostPort(s.config.Server.Host, s.config.Server.Port),
		Handler: s.router,
	}

	// Канал для graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("Server error: %v", err)
		}
	}()

	s.log.Infof("Server is running on host %s and port %s", s.config.Server.Host, s.config.Server.Port)

	// Ожидание сигнала для graceful shutdown
	<-done
	s.log.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		s.log.Fatalf("Server shutdown error: %v", err)
	}

	s.log.Info("Server stopped")
	return nil
}
