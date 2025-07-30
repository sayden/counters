package main

import (
	"context"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/sayden/counters"
	"github.com/sayden/counters/server"
	"github.com/spf13/afero"
)

type VirtualFileSystem struct {
	afero.Fs
}

func (v *VirtualFileSystem) GenerateCounters(ctx context.Context, byt []byte, s ...server.Subscriber) (server.CounterImages, error) {
	var subscriber server.Subscriber = &server.NoOpSubscriber{}
	if len(s) > 0 {
		subscriber = s[0]
	}

	filenamesInUse := &sync.Map{}
	tempTemplate, err := counters.ParseCounterTemplate(byt, filenamesInUse)
	if err != nil {
		return nil, err
	}

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	err = os.Chdir(os.ExpandEnv(tempTemplate.WorkingDirectory))
	if err != nil {
		return nil, err
	}

	newTemplate, err := tempTemplate.ParsePrototype()
	if err != nil {
		return nil, err
	}

	response := server.CounterImages(make([]server.CounterImage, 0, len(newTemplate.Counters)))

	// Generate a bunch of templates.CounterImage in parallel
	// Those templates do not leave the scope of this function
	i := 0
	fileNumberPlaceholder := 0
	var wg sync.WaitGroup
	wg.Add(len(newTemplate.Counters))
	var ch = make(chan server.CounterImage, len(newTemplate.Counters))

	subscriber.Total(len(newTemplate.Counters))

	for _, counter := range newTemplate.Counters {
		go func(counter counters.Counter, i, fileNumberPlaceholder int) {
			defer wg.Done()

			filename := counter.GenerateCounterFilename("", i, filenamesInUse)
			filename = "/" + filename
			file, err := v.Create(filename)
			if err != nil {
				log.Error("Could not create file in VFS", "filename", filename, "error", err)
				return
			}

			if err = counter.EncodeCounter(file, false); err != nil {
				log.Error(err)
				return
			}
			subscriber.OnEvent(ctx, i)

			defer file.Close()

			ch <- server.CounterImage{
				Filename:   "/api/images" + filename,
				Id:         counter.Filename,
				PrettyName: strings.ReplaceAll(counter.PrettyName, "_", " "),
			}
		}(counter, i, fileNumberPlaceholder)

		i++
		fileNumberPlaceholder++
	}

	wg.Wait()

	// Get all the templates.CounterImage that have been generated into an array
	for counterImage := range ch {
		response = append(response, counterImage)
		if len(response) == cap(response) {
			break
		}
	}

	// Sort the array by the counter's filename so that their order is consistent in the browser
	slices.SortFunc(response, func(a, b server.CounterImage) int {
		if a.Id < b.Id {
			return -1
		} else if a.Id > b.Id {
			return 1
		}

		return 0
	})

	return response, nil
}
