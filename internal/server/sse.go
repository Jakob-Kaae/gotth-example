package server

import (
	"fmt"
	"net/http"
)

func (s *server) parkingSSE(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("New SSE client connected: %s\n", r.RemoteAddr)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientChan := make(chan string)
	s.sseClientsMu.Lock()
	s.sseClients[clientChan] = struct{}{}
	s.sseClientsMu.Unlock()

	defer func() {
		s.sseClientsMu.Lock()
		delete(s.sseClients, clientChan)
		s.sseClientsMu.Unlock()
		close(clientChan)
	}()

	// Stream messages to this client
	flusher := w.(http.Flusher)

	for html := range clientChan {
		fmt.Fprintf(w, "event: parkingSpots\ndata: %s\n\n", html)
		flusher.Flush()
	}
}
