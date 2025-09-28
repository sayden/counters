package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"

	"desktop/backend"
)

//go:embed all:documentation/dist
var docsAssets embed.FS

//go:embed all:frontend/dist
var frontendAssets embed.FS

func main() {
	log.SetLevel(log.DebugLevel)

	// Create an instance of the app structure
	vfs := &backend.VirtualFileSystem{Fs: afero.NewMemMapFs()}
	appComm := make(chan []byte)
	app := backend.NewApp(vfs, appComm)

	// Create application with options
	err := wails.Run(
		&options.App{
			Title:  "Counters Visualizer",
			Width:  1400,
			Height: 960,
			// WindowStartState: options.Minimised,
			AssetServer: &assetserver.Options{
				Middleware: middleware,
				Handler:    app.Router.GetServerHandler(),
				Assets:     frontendAssets,
			},
			BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
			OnStartup:        app.Startup,
			Bind:             []any{app},
			OnShutdown:       app.Close,
			// StartHidden:      true,
		})
	if err != nil {
		log.Fatal(err)
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Debug("middleware", "path", r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/documentation/dist") {

			// Remvoe the first slash
			w.Header().Set("Content-Type", "text/html")
			r.URL.Path = strings.Replace(r.URL.Path, "/", "", 1)

			// Remove the last slash
			if r.URL.Path[len(r.URL.Path)-1] == '/' {
				r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
			}

			// Read the resulting file or directory
			byt, err := docsAssets.ReadFile(r.URL.Path)
			if err != nil {
				if pathErr, ok := err.(*fs.PathError); !ok || pathErr.Err.Error() != "is a directory" {
					log.Error("Could not read file", "error", err)
					return
				}

				// Path is actually a directory, so try to read the index.html file
				byt, err = docsAssets.ReadFile(r.URL.Path + "/index.html")
				if err != nil {
					log.Error("Could not read file", "unhandled error", err)
					return
				}
			}

			_, _ = w.Write(byt)

			return
		}

		next.ServeHTTP(w, r)
	})
}
