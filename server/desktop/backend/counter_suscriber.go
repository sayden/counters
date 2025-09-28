package backend

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type CountersSuscriber struct {
	total int
	sync.Mutex
}

func (c *CountersSuscriber) OnEvent(wCtx context.Context, n int) {
	c.Lock()
	defer c.Unlock()
	c.total--
	runtime.EventsEmit(wCtx, "processed_left", c.total)
}

func (c *CountersSuscriber) Total(t int) {
	c.total = t
}
