package main

import (
	"fmt"
	"github.com/uptrace/bunrouter"
	"log"
	"net/http"
	"sparky-back/internal/middlewares"
	"time"
)

func main() {
	router := bunrouter.New(
		bunrouter.Use(middlewares.Log),
	)
	router.POST("/events", func(w http.ResponseWriter, req bunrouter.Request) error {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return nil
		}

		// Set the necessary headers to allow Server-Sent Events
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		for {
			// Write some data to the client
			fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.Stamp))

			// Flush the data immediately instead of buffering it for later.
			flusher.Flush()

			// Pause for a second before the next iteration.
			time.Sleep(time.Second)
		}
	})
	handler := http.HandlerFunc(router.ServeHTTP)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
