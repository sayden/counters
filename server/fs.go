package server

import "context"

type Subscriber interface {
	OnEvent(ctx context.Context, n int)
	Total(int)
}

type Filesystem interface {
	GenerateCounters(ctx context.Context, byt []byte, sub ...Subscriber) (CounterImages, error)
}
