package main

import (
	"context"
	"github.com/Serasmi/home-library/internal/handler"
	"github.com/Serasmi/home-library/internal/repository"
	"github.com/Serasmi/home-library/internal/repository/mongorepo"
	"github.com/Serasmi/home-library/internal/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file: %s", err.Error())
	}

	ctx := context.Background()

	router := mux.NewRouter()

	mongoClient := mongorepo.New()
	if err = mongoClient.Init(ctx); err != nil {
		logrus.Fatalf("Error connecting mongo server: %s", err.Error())
	}
	defer func() {
		if err := mongoClient.Close(ctx); err != nil {
			logrus.Fatalf("Error closing mongo server: %s", err.Error())
		}
	}()

	repos := repository.NewMongoRepository(mongoClient)
	services := service.New(repos)
	handlers := handler.New(services)

	router.Handle(handlers.HealthCheck())
	router.Handle(handlers.GetAllBooks())

	srv := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 0,
	}

	go func() {
		// catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		c, cancel := context.WithTimeout(ctx, 15) // TODO: move wait timout to flag
		defer cancel()

		if err := srv.Shutdown(c); err != nil {
			logrus.Fatal("HTTP server Shutdown: %s", err.Error())
		}

		logrus.Info("Shutting down...")
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logrus.Fatalf("Error listening server: %s", err.Error())
	}
}
