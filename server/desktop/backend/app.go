package backend

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/sayden/counters"
	"github.com/sayden/counters/server"
)

var startingFolder, _ = os.Getwd()

// App struct
type App struct {
	WailsCtx context.Context
	server.Filesystem
	server.Subscriber
	*FileWatcher

	appComm chan []byte

	Router HttpRouter
}

// NewApp creates a new App application struct
func NewApp(fs server.Filesystem, appComm chan []byte) *App {
	app := &App{Filesystem: fs, Subscriber: &CountersSuscriber{}, appComm: appComm}
	app.setupServer(fs, appComm)
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.WailsCtx = ctx
	a.Router.wailsCtx = a.WailsCtx
	go a.listenApiCode()
}

// GetCounters stored in the global store, but does no produce any side effect
func (a *App) GetCounters() []server.CounterImage {
	server.GlobalStore.Lock()
	defer server.GlobalStore.Unlock()

	return server.GlobalStore.CounterImages
}

// SelectFile opens a file dialog and the selected file path is loaded in memory in a new goroutine
// Then the path is returned
func (a *App) SelectFile() (string, error) {
	result, err := runtime.OpenFileDialog(a.WailsCtx, runtime.OpenDialogOptions{
		Title: "Select a file",
		Filters: []runtime.FileFilter{{
			DisplayName: "All Files",
			Pattern:     "*.json",
		}},
	})
	if err != nil {
		return "", err
	}

	go a.loadFileInMemory(result)

	return result, err
}

// GetImage returns the first image in a template
func (a *App) GetImage(data string) ([]byte, error) {
	newTemplate, err := a.getTemplate([]byte(data))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	if err = newTemplate.Counters[0].EncodeCounter(buf, false); err != nil {
		log.Error("[Wails]", "Could not encode counter", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

// GETCounters from the GlobalStore in JSON
func (a *App) GETCounters(c fiber.Ctx) error {
	server.GlobalStore.Lock()
	defer server.GlobalStore.Unlock()

	return c.JSON(fiber.Map{"counters": server.GlobalStore.CounterImages})
}

func (a *App) SSE(c fiber.Ctx) error {
	// Set headers for SSE
	c.Response().Header.Set("Content-Type", "text/event-stream")
	c.Response().Header.Set("Cache-Control", "no-cache")
	c.Response().Header.Set("Connection", "keep-alive")
	c.Response().Header.Set("Access-Control-Allow-Origin", "*")

	// Create a channel to notify of client disconnect
	clientChan := make(chan bool)
	go func() {
		<-c.Done()
		clientChan <- true
	}()

	flusher, ok := c.Response().BodyWriter().(http.Flusher)
	if !ok {
		return errors.New("writer does not support flushing")
	}

	// Send events
	for {
		select {
		case <-clientChan:
			log.Debug("Client disconnected")
			return nil

		// received in /api/code endpoint
		case schemaByt := <-a.appComm:
			counters, err := a.GenerateCounters(a.WailsCtx, schemaByt)
			if err != nil {
				log.Error(err)
				continue
			}

			server.GlobalStore.Lock()
			server.GlobalStore.CounterImages = counters
			server.GlobalStore.Unlock()

			// Signal the frontend that it can fetch the new data using /state endpoint
			// which returns just the grid of counter to use in HTMX
			_, err = fmt.Fprintf(c.Response().BodyWriter(), "event: Grid\ndata:ok\n\n")
			if err != nil {
				log.Error("Could not write to the client", "error", err)
			}
			flusher.Flush()
		}
	}
}

func (a *App) Close(ctx context.Context) {
	a.FileWatcher.Close()
}

func (a *App) listenApiCode() {
	var err error

	for byt := range a.appComm {
		server.GlobalStore.Lock()
		server.GlobalStore.CounterImages, err = a.GenerateCounters(a.WailsCtx, byt, a.Subscriber)
		if err != nil {
			log.Error("Could not generate counters", "error", err)
			err = nil
		} else {
			// Notify the frontend that the counters have changed
			runtime.EventsEmit(a.WailsCtx, "counters")
		}
		server.GlobalStore.Unlock()
	}

	log.Warn("AppComm channel closed")
}

func (a *App) getTemplate(data []byte) (*counters.CounterTemplate, error) {
	filenamesInUse := &sync.Map{}
	template, err := counters.ParseCounterTemplate(data, filenamesInUse)
	if err != nil {
		log.Error("[Wails]", "could not parse template", err)
		return nil, err
	}

	if err = os.Chdir(os.ExpandEnv(template.WorkingDirectory)); err != nil {
		log.Error("[Wails]", "Could not change working directory", err)
		return nil, errors.New("could not change working directory")
	}

	newTemplate, err := template.ParsePrototype()
	if err != nil {
		log.Error("[Wails]", "Could not parse prototype", err)
		return nil, err
	}

	if len(newTemplate.Counters) == 0 {
		log.Error("[Wails]", "No counters found in template")
		return nil, errors.New("no counters found in template")
	}

	return newTemplate, nil
}

func (a *App) loadFileInMemory(filename string) {
	byt, err := os.ReadFile(filename)
	if err != nil {
		log.Error("Could not read file", "error", err)
	}

	runtime.EventsEmit(a.WailsCtx, "processed_left")

	server.GlobalStore.Lock()
	server.GlobalStore.CounterImages, err = a.GenerateCounters(a.WailsCtx, byt, a.Subscriber)
	server.GlobalStore.Unlock()
	if err != nil {
		log.Error("Could not generate counters", "error", err)
		return
	}

	// Notify the frontend that the counters have changed
	runtime.EventsEmit(a.WailsCtx, "counters")

	// Close current watcher if any
	if a.FileWatcher != nil {
		a.FileWatcher.Close()
	}

	a.FileWatcher, err = NewFileWatcher(filename, &a.Router)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Watch(a.WailsCtx)
	if err != nil {
		log.Error("Could not watch file", "error", err)
	}
}

func (a *App) setupServer(vfs server.Filesystem, appComm chan []byte) {
	router := fiber.New()

	// App routes
	router.Get("/api/images/:filename", a.Router.get)
	router.Get("*", func(c fiber.Ctx) error {
		log.Debug("FALLBACK", "method", "GET", "path", c.Path())
		return nil
	})

	a.Router.fiberServeHTTP = adaptor.FiberApp(router)
	a.Router.wailsCtx = a.WailsCtx
	a.Router.Filesystem = vfs
	a.Router.appComm = appComm
}
