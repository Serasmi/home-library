package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Serasmi/home-library/pkg/files"

	"github.com/Serasmi/home-library/internal/api/books"
	"github.com/Serasmi/home-library/internal/api/health"
	"github.com/Serasmi/home-library/internal/api/uploads"
	"github.com/Serasmi/home-library/internal/auth"
	"github.com/Serasmi/home-library/internal/config"
	apiRouter "github.com/Serasmi/home-library/internal/router"
	"github.com/Serasmi/home-library/internal/user"
	"github.com/Serasmi/home-library/pkg/logging"
	"github.com/Serasmi/home-library/pkg/mongodb"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

const (
	apiPath           = "/api"
	booksCollection   = "books"
	uploadsCollection = "uploads"
	usersCollection   = "users"
)

func main() {
	logger := logging.NewLogger()

	cfg := config.InitConfig(logger)

	logger.SetLevel(cfg.App.LogLevel)

	ctx := context.Background()
	router := apiRouter.NewRouter(logger)

	mongoClient, err := mongodb.NewClient(ctx, cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.Name)
	if err != nil {
		logger.Fatalf("Error connecting mongodb server: %s", err.Error())
	}

	defer func() {
		// TODO: make mongoClient closable and call Close method.
		//  Close method should use context.WithTimeout().
		if err := mongoClient.Disconnect(ctx); err != nil {
			logger.Fatalf("Error closing mongodb server: %s", err.Error())
		}
	}()

	fileProvider := files.NewFSProvider(logger)

	userStorage := user.NewMongoStorage(mongoClient.Database(cfg.DB.Name), usersCollection, logger)
	userService := user.NewService(userStorage, logger)

	authHandler := auth.NewHandler(userService, logger)
	authHandler.Register(router)

	healthHandler := health.NewHandler(apiPath)
	healthHandler.Register(router)

	uploadsStorage := uploads.NewMongoStorage(mongoClient.Database(cfg.DB.Name), uploadsCollection, logger)
	uploadsUseCase := uploads.NewUseCase(uploadsStorage, fileProvider, logger)
	uploadsHandler := uploads.NewHandler(apiPath, uploadsUseCase, logger)
	uploadsHandler.Register(router)

	booksStorage := books.NewMongoStorage(mongoClient.Database(cfg.DB.Name), booksCollection, logger)
	// booksStorage := db.NewMockStorage(logger)
	booksUseCase := books.NewUseCase(booksStorage, uploadsStorage, fileProvider, logger)
	booksHandler := books.NewHandler(apiPath, booksUseCase, logger)
	booksHandler.Register(router)

	start(ctx, router, logger, cfg)
}

func start(ctx context.Context, router *httprouter.Router, logger *logging.Logger, cfg *config.Config) {
	logger.Infof("Start application")

	listener, err := net.Listen("tcp", cfg.App.Host+":"+cfg.App.Port)
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		if err := srv.Serve(listener); err != http.ErrServerClosed {
			logger.Fatalf("Error listening server: %s", err.Error())
		}
	}()

	logger.Infof("Server is listening on %s:%s", cfg.App.Host, cfg.App.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	c, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(c); err != nil {
		logrus.Fatalf("HTTP server Shutdown: %s", err.Error())
	}

	logrus.Info("Shutting down...")
}
