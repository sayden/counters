package httphandlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/sayden/counters/server"
	templates "github.com/sayden/counters/server/templates/generated"
)

func NewWeb() *Web {
	return &Web{}
}

type Web struct{}

// GETGrid of counters to use in HTMX
func (web *Web) GETGrid(c *fiber.Ctx) error {
	server.GlobalStore.Lock()
	defer server.GlobalStore.Unlock()

	c.Response().Header.Set("Cache-Control", "no-cache")

	return templates.Counters(server.GlobalStore.CounterImages).
		Render(c.Context(), c.Request().BodyWriter())
}

func (web *Web) GETIndex(c *fiber.Ctx) error {
	server.GlobalStore.Lock()
	defer server.GlobalStore.Unlock()

	c.Response().Header.Set("Cache-Control", "no-cache")

	return templates.Index().
		Render(c.Context(), c.Request().BodyWriter())
}
