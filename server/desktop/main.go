package main

import (
	"bytes"
	"embed"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/spf13/afero"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/sayden/counters"
)

var startingFolder, _ = os.Getwd()

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	log.SetLevel(log.DebugLevel)

	// Create an instance of the app structure
	vfs := &VirtualFileSystem{Fs: afero.NewMemMapFs()}
	app := NewApp(vfs, &countersSuscriber{})

	// Create application with options
	if err := wails.Run(&options.App{
		Title:            "Counters Visualizer",
		Width:            1400,
		Height:           960,
		WindowStartState: options.Maximised,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: setupServer(vfs),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind:             []any{app},
		OnShutdown:       app.Close,
	}); err != nil {
		log.Fatal(err)
	}
}

func setupServer(vfs *VirtualFileSystem) http.Handler {
	server := fiber.New()
	hh := &httpHandler{VirtualFileSystem: vfs}

	// Setup routes
	server.Post("/temp.png", hh.New)
	server.Get("/api/images/:filename", hh.Get)

	adaptor.HTTPHandler(hh)
	handler := adaptor.FiberApp(server)
	hh.fiberHandler = handler

	return handler
}

type httpHandler struct {
	*VirtualFileSystem
	fiberHandler http.HandlerFunc
}

func (h *httpHandler) New(c *fiber.Ctx) error {
	byt := c.Body()
	filenamesInUse := &sync.Map{}
	template, err := counters.ParseCounterTemplate(byt, filenamesInUse)
	if err != nil {
		log.Error("Could not parse template", err)
		return err
	}

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	if err = os.Chdir(os.ExpandEnv(template.WorkingDirectory)); err != nil {
		log.Error("Could not change working directory", "error", err)
		return err
	}

	newTemplate, err := template.ParsePrototype()
	if err != nil {
		log.Error("Could not parse prototype", "error", err)
		return err
	}

	if len(newTemplate.Counters) == 0 {
		log.Error("No counters found in template")
		return errors.New("no counters found in template")
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	if err = newTemplate.Counters[0].EncodeCounter(buf, false); err != nil {
		log.Error("Could not encode counter", "error", err)
		return err
	}

	c.Response().Header.Set("Content-Type", "image/png")
	n, err := c.Write(buf.Bytes())
	// n, err := c.WriteString("data:image/png;base64," + buf.String())
	if err != nil {
		log.Error("Could not write to response", "error", err)
		return err
	}
	log.Debug("Wrote", "bytes", n)

	return nil
}

func (h *httpHandler) Get(c *fiber.Ctx) error {
	// clean the RequestURI to extract the file name
	filename := strings.ReplaceAll("/"+c.Params("filename"), "%20", " ")
	f, err := h.Open(filename)
	if err != nil {
		log.Error("Could not open file", "filename", "/"+c.Params("filename"), "error", err)
		return err
	}
	defer f.Close()

	byt, err := io.ReadAll(f)
	if err != nil {
		log.Error("Could not read file", "filename", "/"+c.Params("filename"), "error", err)
		return err
	}

	if _, err = c.Write(byt); err != nil {
		log.Error("error writing contents of file into response body",
			"filename", "/"+c.Params("filename"), "error", err)
		return err
	}

	return nil
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fiberHandler(w, r)
}
