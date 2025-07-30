package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"os"
	"slices"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/sayden/counters"
)

type CounterImage struct {
	Filename     string `json:"filename"`
	CounterImage string `json:"counter"`
	Id           string `json:"id"`
	PrettyName   string `json:"pretty_name"`
}

type CounterImages []CounterImage

type ResponseMutex struct {
	sync.Mutex
	CounterImages
}

var GlobalStore ResponseMutex

type Base64ImagesFs struct {
}

type NoOpSubscriber struct{}

func (n *NoOpSubscriber) OnEvent(_ context.Context, _ int) {}
func (n *NoOpSubscriber) Total(t int)                      {}

// GenerateCounters in parallel
func (b *Base64ImagesFs) GenerateCounters(byt []byte, s ...Subscriber) (CounterImages, error) {
	var subscriber Subscriber = &NoOpSubscriber{}
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
	os.Chdir(os.ExpandEnv(tempTemplate.WorkingDirectory))

	newTemplate, err := tempTemplate.ParsePrototype()
	if err != nil {
		return nil, err
	}

	response := CounterImages(make([]CounterImage, 0, len(newTemplate.Counters)))

	// Generate a bunch of templates.CounterImage in parallel
	// Those templates do not leave the scope of this function
	i := 0
	fileNumberPlaceholder := 0
	var wg sync.WaitGroup
	wg.Add(len(newTemplate.Counters))
	var ch = make(chan CounterImage, len(newTemplate.Counters))

	subscriber.Total(len(newTemplate.Counters))

	for _, counter := range newTemplate.Counters {
		go func(counter counters.Counter, i, fileNumberPlaceholder int) {
			defer wg.Done()
			buf := new(bytes.Buffer)
			wc := base64.NewEncoder(base64.StdEncoding, buf)

			// get a canvas with the rendered counter. The canvas can be written to a io.Writer
			err := counter.EncodeCounter(wc, newTemplate.DrawGuides)
			if err != nil {
				log.Error(err)
				return
			}

			subscriber.OnEvent(context.Background(), i)

			counter.GenerateCounterFilename("", i, filenamesInUse)
			counterImage := CounterImage{
				CounterImage: "data:image/png;base64," + buf.String(),
				Id:           counter.Filename,
				PrettyName:   counter.PrettyName,
			}
			wc.Close()
			ch <- counterImage
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
	slices.SortFunc(response, func(a, b CounterImage) int {
		if a.Id < b.Id {
			return -1
		} else if a.Id > b.Id {
			return 1
		}

		return 0
	})

	return response, nil
}
