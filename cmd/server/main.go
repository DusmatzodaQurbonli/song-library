package main

import (
	"context"
	"fmt"
	"github.com/DusmatzodaQurbonli/song-library/internal/config"
	"github.com/DusmatzodaQurbonli/song-library/internal/handler"
	log "github.com/DusmatzodaQurbonli/song-library/internal/logger"
	"github.com/DusmatzodaQurbonli/song-library/internal/repository"
	"github.com/DusmatzodaQurbonli/song-library/internal/service"
	"github.com/DusmatzodaQurbonli/song-library/pkg/http"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	fx.New(
		fx.Provide(
			config.New,
			log.New,
			initializeDB,
			repository.NewSongRepository,
			service.NewMusicInfoClient, // Теперь передаем правильно
			service.NewSongService,
			handler.NewSongHandler,
			http.NewServer,
		),
		fx.Invoke(runMigrations),
		fx.Invoke(startServer),
	).Run()
}

func initializeDB(cfg *config.Config, log *logrus.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Pass, cfg.DB.Name, cfg.DB.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	return db, nil
}

func runMigrations(db *gorm.DB, log *logrus.Logger, cfg *config.Config) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}
	defer sqlDB.Close()

	m, err := migrate.New(
		fmt.Sprintf("file://%s", "internal/migration/migrations"),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name,
		),
	)
	if err != nil {
		log.Fatal("Migration initialization failed: ", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Info("No new migrations to apply.")
		} else {
			log.Fatal("Failed to run migrations: ", err)
		}
	} else {
		log.Info("Migrations applied successfully!")
	}
}

func startServer(lc fx.Lifecycle, srv *http.Server, cfg *config.Config, log *logrus.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Infof("Starting server on port %s", cfg.Server.Port)
				if err := srv.Run(); err != nil {
					log.Fatal("Failed to start server: ", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down server...")
			return nil
		},
	})
}
