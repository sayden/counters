package main

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type countersSuscriber struct {
	total int
	sync.Mutex
}

func (c *countersSuscriber) OnEvent(ctx context.Context, n int) {
	c.Lock()
	defer c.Unlock()
	c.total--
	runtime.EventsEmit(ctx, "processed_left", c.total)
}

func (c *countersSuscriber) Total(t int) {
	c.total = t
}
