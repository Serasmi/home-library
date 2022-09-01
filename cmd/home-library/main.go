package main

import (
	"context"
	"github.com/Serasmi/home-library/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file: %s", err.Error())
	}

	ctx := context.Background()

	srv := service.New()

	srv.Run(ctx)
}
