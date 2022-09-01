package service

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s Service) Run(ctx context.Context) {
	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Errorf("Error encoding healthcheck response: %s", err.Error())
		}
	})

	srv := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 0,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logrus.Fatalf("Error listening server: %s", err.Error())
		}
	}()

	// catch signal and invoke graceful termination
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(ctx, 15) // TODO: move wait timout to flag
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("HTTP server Shutdown: %s", err.Error())
	}

	logrus.Info("Shutting down...")

	os.Exit(0)
}
