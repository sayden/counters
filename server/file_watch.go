package server

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
	fileToWatch string
	watcher     *fsnotify.Watcher

	listeners []FileWatchListener
}

func NewFileWatcher(filepath string, listeners ...FileWatchListener) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FileWatcher{
		fileToWatch: filepath,
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
				log.Debug("Event", "name", event.Name, "op", event.Op)
				if !ok {
					continue
				}
				if path.Ext(event.Name) != ".json" {
					continue
				}

				if event.Op == fsnotify.Write {
					for _, listener := range fw.listeners {
						// Wait for the file to be closed
						var info os.FileInfo
						var err error
						for range 10 {
							if info, err = os.Stat(event.Name); err == nil && info.Size() > 0 {
								if err = listener.OnEvent(ctx, event.Op.String(), event.Name); err != nil {
									log.Error("OnEvent listener", "error", err)
								}
								break
							}
							time.Sleep(200 * time.Millisecond)
						}
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
	return fw.watcher.Add(fw.fileToWatch)
}
