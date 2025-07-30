package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/http2"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/sayden/counters/server"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var startingFolder, _ = os.Getwd()

func main() {
	log.SetLevel(log.DebugLevel)

	var fileToWatch string

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Launch the server",
	}
	cmd.Flags().StringVarP(&fileToWatch, "file", "f", "", "File/Folder to watch")

	ctx := context.Background()
	err := fang.Execute(ctx, cmd)
	if err != nil {
		os.Exit(1)
	}

	router := gin.Default()

	router.Static("/static", "cmd/static")

	// API routes
	api := NewApiHandler(
		&server.VirtualFileSystem{Fs: afero.NewMemMapFs()},
		make(chan []byte))
	router.POST("/api/code", api.POSTCode)
	router.GET("/api/sse", api.SSE)
	router.GET("/api/counters", api.GETCounters)

	// Watch file if provided
	var fw *server.FileWatcher
	if fileToWatch != "" {
		if fw, err = server.NewFileWatcher(fileToWatch, api); err != nil {
			log.Fatal(err)
		}
		defer fw.Close()

		if err = fw.Watch(ctx); err != nil {
			log.Fatal(err)
		}
	}

	// Capture SIGINT and SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(fw *server.FileWatcher) {
		<-c
		log.Debug("Caught signal, shutting down...")
		if err := fw.Close(); err != nil {
			log.Error(err)
		}
		log.Debug("Server closed")
		os.Exit(0)
	}(fw)

	// Web routes
	web := webHandler{}
	router.GET("/state", web.GETGrid)
	router.GET("/", web.GETIndex)

	// Launch the server
	server := &http.Server{Addr: ":8090", Handler: router}
	log.Info("Server is running on http://localhost:8090")
	log.Info("This server must be run in the same folder as the template to find its assets")

	// Configure HTTP/2
	if err := http2.ConfigureServer(server, &http2.Server{}); err != nil {
		log.Fatal(err)
	}

	log.Error(server.ListenAndServe())
}
