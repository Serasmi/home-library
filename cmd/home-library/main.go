package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file: %s", err.Error())
	}

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

	logrus.Fatal(srv.ListenAndServe())
}
