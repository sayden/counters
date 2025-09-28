package backend

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v3"

	"github.com/sayden/counters"
	"github.com/sayden/counters/server"
)

type HttpRouter struct {
	server.Filesystem

	wailsCtx       context.Context
	fiberServeHTTP http.HandlerFunc
	appComm        chan []byte
}

func (a *HttpRouter) GetServerHandler() http.Handler {
	return a.fiberServeHTTP
}

func (a *HttpRouter) get(c fiber.Ctx) error {
	// clean the RequestURI to extract the file name
	filename := strings.ReplaceAll("/"+c.Params("filename"), "%20", " ")

	vfs, ok := a.Filesystem.(*VirtualFileSystem)
	if !ok {
		return errors.New("filesystem is not a VirtualFileSystem")
	}

	f, err := vfs.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	byt, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if _, err = c.Write(byt); err != nil {
		return err
	}

	return nil
}

func (a *HttpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.fiberServeHTTP(w, r)
}

// OnEvent is the implementation of fileWatchListener
func (a *HttpRouter) OnEvent(ctx context.Context, event, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	byt, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return a.readReadCloser(ctx, byt)
}

func (a *HttpRouter) readReadCloser(ctx context.Context, byt []byte) error {
	if err := counters.ValidateSchemaBytes[counters.CounterTemplate](byt); err != nil {
		return err
	}

	log.Info("Template received", "length", len(byt))

	a.appComm <- byt

	log.Info("Template sent", "length", len(byt))

	return nil
}
