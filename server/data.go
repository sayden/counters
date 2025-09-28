package server

import (
	"context"
	"sync"
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

type NoOpSubscriber struct{}

func (n *NoOpSubscriber) OnEvent(_ context.Context, _ int) {}
func (n *NoOpSubscriber) Total(t int)                      {}
