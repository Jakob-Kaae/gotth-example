package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Jakob-Kaae/gotth-example/internal/parking/infrastructure/httpclient"
	"github.com/Jakob-Kaae/gotth-example/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "[Spooktober] ", log.LstdFlags)

	port := 9000

	logger.Print("Creating guests store..")

	httpClient := &http.Client{Timeout: 5 * time.Second}
	apiClient := httpclient.NewParkingAPIClient("https://letparkeringapi.azurewebsites.net", httpClient)

	srv, err := server.NewServer(logger, port, apiClient)
	if err != nil {
		logger.Fatalf("Error when creating server: %s", err)
		os.Exit(1)
	}
	if err := srv.Start(); err != nil {
		logger.Fatalf("Error when starting server: %s", err)
		os.Exit(1)
	}
}
