package main

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
)

type FileWatchListener interface {
	OnEvent(ctx context.Context, event, filepath string) error
}

type FileWatcher struct {
	fileToWatch, folder string
	watcher             *fsnotify.Watcher

	listeners []FileWatchListener
}

func NewFileWatcher(filepath string, listeners ...FileWatchListener) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	folder := path.Dir(filepath)

	return &FileWatcher{
		fileToWatch: filepath,
		folder:      folder,
		watcher:     watcher,
		listeners:   listeners,
	}, nil
}

func (fw *FileWatcher) Close() error {
	return fw.watcher.Close()
}

func (fw *FileWatcher) Watch(ctx context.Context) error {
	go func() {
		for {
			select {
			case event, ok := <-fw.watcher.Events:
				if !ok {
					continue
				}

				if event.Name != fw.fileToWatch {
					log.Warn("Ignoring event by its filename", "name", event.Name, "op", event.Op)
					continue
				}

				if path.Ext(event.Name) != ".json" {
					continue
				}

				if event.Op != fsnotify.Write {
					log.Warn("Ignoring event by its operation", "name", event.Name, "op", event.Op)
					continue
				}

				var info os.FileInfo
				var err error
				for range 10 {
					// Wait for the file to be closed
					if info, err = os.Stat(event.Name); err == nil && info.Size() > 0 {
						break
					}
					log.Debug("Waiting for file to be closed", "path", event.Name)
					time.Sleep(500 * time.Millisecond)
				}

				for _, listener := range fw.listeners {
					log.Info("Valid Event", "name", event.Name, "op", event.Op)
					if err = listener.OnEvent(ctx, event.Op.String(), event.Name); err != nil {
						log.Error("OnEvent listener", "error", err)
					}
				}

			case err, ok := <-fw.watcher.Errors:
				if !ok {
					return
				}
				log.Error("ERROR:", err)
			}
		}
	}()

	// Add a file or directory to watch
	return fw.watcher.Add(fw.folder)
}
