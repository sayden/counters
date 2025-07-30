package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/sayden/counters"
	"github.com/sayden/counters/server"
)

func NewApiHandler(fs server.Filesystem, codeCh chan []byte) *apiHandler {
	return &apiHandler{fs: fs, codeCh: codeCh}
}

type apiHandler struct {
	codeCh chan []byte
	fs     server.Filesystem
}

// OnEvent is the implementation of fileWatchListener
func (api *apiHandler) OnEvent(ctx context.Context, event, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	return api.readReadCloser(ctx, file)
}

// GETCounters from the GlobalStore in JSON
func (api *apiHandler) GETCounters(c *gin.Context) {
	server.GlobalStore.Lock()
	defer server.GlobalStore.Unlock()

	c.JSON(http.StatusOK, gin.H{"counters": server.GlobalStore.CounterImages})
}

// SSE is an endpoint that the frontend will connect to to receive updates via SSE
func (api *apiHandler) SSE(c *gin.Context) {
	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	if c.Request.ProtoMajor == 2 {
		log.Debug("Client is using HTTP/2")
	} else {
		log.Debug("Client is using HTTP/1.x")
	}

	// Create a channel to notify of client disconnect
	clientChan := make(chan bool)
	go func() {
		<-c.Request.Context().Done()
		clientChan <- true
	}()

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		http.Error(c.Writer, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send events
	for {
		select {
		case <-clientChan:
			log.Debug("Client disconnected")
			return

		// received in /api/code endpoint
		case schemaByt := <-api.codeCh:
			counters, err := api.fs.GenerateCounters(schemaByt)
			if err != nil {
				log.Error(err)
				continue
			}

			server.GlobalStore.Lock()
			server.GlobalStore.CounterImages = counters
			server.GlobalStore.Unlock()

			// Signal the frontend that it can fetch the new data using /state endpoint
			// which returns just the grid of counter to use in HTMX
			_, err = fmt.Fprintf(c.Writer, "event: Grid\ndata:ok\n\n")
			if err != nil {
				log.Error("Could not write to the client", "error", err)
			}
			flusher.Flush()
		}
	}
}

// POSTCode here to send it to the frontend
func (api *apiHandler) POSTCode(c *gin.Context) {
	defer func() {
		if err := os.Chdir(startingFolder); err != nil {
			log.Error(err)
		}
	}()

	err := api.readReadCloser(c.Request.Context(), c.Request.Body)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (api *apiHandler) readReadCloser(ctx context.Context, reader io.ReadCloser) error {
	defer reader.Close()

	byt, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if err = counters.ValidateSchemaBytes[counters.CounterTemplate](byt); err != nil {
		return err
	}

	api.codeCh <- byt

	return nil
}
