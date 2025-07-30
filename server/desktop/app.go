package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/charmbracelet/log"

	"github.com/sayden/counters"
	"github.com/sayden/counters/server"
)

// App struct
type App struct {
	ctx context.Context

	*FileWatcher
	server.Filesystem
	server.Subscriber
}

// NewApp creates a new App application struct
func NewApp(fs server.Filesystem, suscriber server.Subscriber) *App {
	return &App{Filesystem: fs, Subscriber: suscriber}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetCounters() []server.CounterImage {
	server.GlobalStore.Lock()
	defer server.GlobalStore.Unlock()

	return server.GlobalStore.CounterImages
}

// This function opens a file dialog and returns the selected file path
func (a *App) SelectFile() (string, error) {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a file",
		Filters: []runtime.FileFilter{{
			DisplayName: "All Files",
			Pattern:     "*.json",
		}},
	})
	if err != nil {
		return "", err
	}

	log.Debug("Selected file", "path", result)

	go a.loadFileInMemory(result)

	return result, err
}

func (a *App) loadFileInMemory(filename string) ([]byte, error) {
	byt, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	runtime.EventsEmit(a.ctx, "processed_left")

	server.GlobalStore.Lock()
	server.GlobalStore.CounterImages, err = a.GenerateCounters(a.ctx, byt, a.Subscriber)
	server.GlobalStore.Unlock()
	if err != nil {
		return nil, err
	}

	// Notify the frontend that the counters have changed
	runtime.EventsEmit(a.ctx, "counters")

	// Close current watcher if any
	if a.FileWatcher != nil {
		a.FileWatcher.Close()
	}

	a.FileWatcher, err = NewFileWatcher(filename, a)
	if err != nil {
		log.Fatal(err)
	}

	return nil, a.Watch(a.ctx)
}

func (a *App) OnEvent(ctx context.Context, event, filepath string) error {
	log.Info("Change detected", "path", filepath)

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	byt, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if err = counters.ValidateSchemaBytes[counters.CounterTemplate](byt); err != nil {
		return err
	}

	server.GlobalStore.Lock()
	server.GlobalStore.CounterImages, err = a.GenerateCounters(a.ctx, byt, a.Subscriber)
	server.GlobalStore.Unlock()
	if err != nil {
		return err
	}

	log.Debug("Emitting counters", "count", len(server.GlobalStore.CounterImages))

	runtime.EventsEmit(a.ctx, "counters")

	return err
}

func (a *App) GetImage(data string) ([]byte, error) {
	filenamesInUse := &sync.Map{}
	template, err := counters.ParseCounterTemplate([]byte(data), filenamesInUse)
	if err != nil {
		log.Error("Could not parse template", err)
		return nil, err
	}

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	if err = os.Chdir(os.ExpandEnv(template.WorkingDirectory)); err != nil {
		log.Error("Could not change working directory", "error", err)
		return nil, err
	}

	newTemplate, err := template.ParsePrototype()
	if err != nil {
		log.Error("Could not parse prototype", "error", err)
		return nil, err
	}

	if len(newTemplate.Counters) == 0 {
		log.Error("No counters found in template")
		return nil, errors.New("no counters found in template")
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	if err = newTemplate.Counters[0].EncodeCounter(buf, false); err != nil {
		log.Error("Could not encode counter", "error", err)
		return nil, err
	}

	// return []byte("data:image/png;base64," + buf.String()), nil
	return buf.Bytes(), nil
}

func (a *App) Close(ctx context.Context) {
	log.Info("Closing watchers")
	a.FileWatcher.Close()
}
