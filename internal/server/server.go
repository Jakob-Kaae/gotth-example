package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Jakob-Kaae/gotth-example/internal/parking/application/service"
	"github.com/Jakob-Kaae/gotth-example/internal/parking/domain/model"
	"github.com/Jakob-Kaae/gotth-example/internal/parking/infrastructure/httpclient"
	parkingStore "github.com/Jakob-Kaae/gotth-example/internal/parking/infrastructure/store"
	"github.com/Jakob-Kaae/gotth-example/internal/parking/interface/httpinterface"
	"github.com/Jakob-Kaae/gotth-example/internal/templates"
)

type server struct {
	logger       *log.Logger
	port         int
	httpServer   *http.Server
	parkingStore *parkingStore.InMemoryParkingStore
	apiClient    *httpclient.ParkingAPIClient
	sseClients   map[chan string]struct{}
	sseClientsMu sync.Mutex
}

// Creat a new server instance with the given logger and port
func NewServer(logger *log.Logger, port int, apiClient *httpclient.ParkingAPIClient) (*server, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &server{
		logger:       logger,
		port:         port,
		parkingStore: parkingStore.NewInMemoryParkingStore(),
		apiClient:    apiClient,
		sseClients:   make(map[chan string]struct{}),
	}, nil
}

// Start the server
func (s *server) Start() error {
	s.logger.Printf("Starting server on port %d", s.port)
	var stopChan chan os.Signal

	// define router
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	router.HandleFunc("GET /", s.defaultHandler)
	router.HandleFunc("GET /about", s.aboutHandler)
	router.HandleFunc("GET /health", s.healthCheckHandler)
	// Interface
	handler := httpinterface.NewParkingHandler(s.parkingStore)
	router.HandleFunc("GET /parking", handler.GetParkingSpots)
	router.HandleFunc("GET /events/parking", s.parkingSSE)
	s.parkingStore.OnChange(func(spots []model.ParkingSpot) {
		s.broadcastParkingUpdate(spots)
	})

	// Application
	pollService := service.NewPollParkingSpotsService(s.apiClient, s.parkingStore)

	// Start polling loop
	go func() {
		for {
			err := pollService.Poll()
			if err != nil {
				s.logger.Printf("Error polling parking spots: %v", err)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	// define server
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: router}

	// create channel to listen for signals
	stopChan = make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error when running server: %s", err)
		}
	}()

	<-stopChan

	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error when shutting down server: %v", err)
		return err
	}
	return nil
}

// GET /
func (s *server) defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	parkingSpots := s.parkingStore.GetAll()
	homeTemplate := templates.Home(parkingSpots)
	err := templates.Layout(homeTemplate, "Parking Tracker", "/").Render(r.Context(), w)
	if err != nil {
		s.logger.Printf("Error when rendering home: %v", err)
	}
}

func (s *server) aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	aboutTemplate := templates.About()
	err := templates.Layout(aboutTemplate, "About", "/about").Render(r.Context(), w)
	if err != nil {
		s.logger.Printf("Error when rendering about: %v", err)
	}
}

func (s *server) broadcastParkingUpdate(spots []model.ParkingSpot) {
	// Render Templ component to HTML
	var buf bytes.Buffer
	err := templates.ParkingList(spots).Render(context.Background(), &buf)
	if err != nil {
		s.logger.Printf("failed to render parking list: %v", err)
		return
	}
	html := buf.String()

	// Send to all SSE clients
	s.sseClientsMu.Lock()
	defer s.sseClientsMu.Unlock()

	for ch := range s.sseClients {
		select {
		case ch <- html:
		default:
			// slow client, drop update
		}
	}

}

// GET /health - HealthCheckHandler is a simple handler to check the health of the server
func (s *server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
