package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pavr1/auth/config"
	"github.com/pavr1/auth/handler"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting AuthServer...")

	log := setupLogger()
	config, err := config.NewConfig(log)
	if err != nil {
		log.WithError(err).Error("Error loading config")
		return
	}

	router := http.NewServeMux()
	//need to add this secret to a secret so it can be retrieved from there
	secret := "test"
	authHandler := handler.NewHandler(log, []byte(secret))

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.HandleFunc("/auth/token", authHandler.ServeHTTP)

	log.WithField("port", config.Server.Port).Info("Listening to AuthServer...")
	// Start the HTTP server
	log.Error(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router))
}

func setupLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	logger.SetReportCaller(true)
	logger.SetLevel(log.DebugLevel)

	// Set the output to stdout
	logger.SetOutput(os.Stdout)

	return logger
}
