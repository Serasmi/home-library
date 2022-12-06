package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Serasmi/home-library/internal/api/books"
	"github.com/Serasmi/home-library/internal/api/books/db"
	apiRouter "github.com/Serasmi/home-library/internal/router"
	"github.com/Serasmi/home-library/pkg/logging"
	"github.com/Serasmi/home-library/pkg/mongodb"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

const apiPath = "/api"

func main() {
	logger := logging.GetLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %s", err.Error())
	}

	ctx := context.Background()
	router := apiRouter.NewRouter(logger)

	mongoClient, err := mongodb.NewClient(ctx, "localhost", "27017", "admin", "admin", "HomeLibrary")
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

	booksStorage := db.NewMongoStorage(mongoClient.Database("HomeLibrary"), "books", logger)
	// booksStorage := db.NewMockStorage(logger)
	booksService := books.NewService(booksStorage, logger)
	booksHandler := books.NewHandler(apiPath, booksService, logger)
	booksHandler.Register(router)

	start(ctx, router, logger)
}

func start(ctx context.Context, router *httprouter.Router, logger logging.Logger) {
	port := os.Getenv("PORT")

	logger.Infof("Start application")

	listener, err := net.Listen("tcp", ":"+port)
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

	logger.Infof("Server is listening on 0.0.0.0:%s", port)

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
