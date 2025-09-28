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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/spf13/cobra"

	"github.com/sayden/counters/server"
	"github.com/sayden/counters/server/httphandlers"
)

func main() {
	log.SetLevel(log.DebugLevel)

	var fileToWatch string

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Launch the server",
	}
	cmd.Flags().StringVarP(&fileToWatch, "file", "f", "", "File/Folder to watch")

	err := fang.Execute(context.Background(), cmd)
	if err != nil {
		os.Exit(1)
	}

	// router := gin.Default()
	router := fiber.New()

	router.Static("/static", "cmd/static")

	// API routes
	api := httphandlers.NewApi(
		context.Background(),
		// &backend.VirtualFileSystem{Fs: afero.NewMemMapFs()},
		&server.Base64ImagesFs{},
		make(chan []byte))
	router.Post("/api/code", api.POSTCode)
	router.Get("/api/sse", api.SSE)
	router.Get("/api/counters", api.GETCounters)

	// Watch file if provided
	var fw *server.FileWatcher
	if fileToWatch != "" {
		if fw, err = server.NewFileWatcher(fileToWatch, api); err != nil {
			log.Fatal(err)
		}
		defer fw.Close()

		if err = fw.Watch(context.Background()); err != nil {
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
	web := httphandlers.Web{}
	router.Get("/state", web.GETGrid)
	router.Get("/", web.GETIndex)

	// Launch the server
	server := &http.Server{Addr: ":8090", Handler: adaptor.FiberApp(router)}
	log.Info("Server is running on http://localhost:8090")
	log.Info("This server must be run in the same folder as the template to find its assets")

	// Configure HTTP/2
	if err := http2.ConfigureServer(server, &http2.Server{}); err != nil {
		log.Fatal(err)
	}

	log.Error(server.ListenAndServe())
}
