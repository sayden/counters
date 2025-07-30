package main

import (
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/sayden/counters/server"
	"github.com/sayden/counters/server/templates/generated"
)

type webHandler struct{}

// GETGrid of counters to use in HTMX
func (web *webHandler) GETGrid(c *gin.Context) {
	server.GlobalStore.Lock()
	defer server.GlobalStore.Unlock()

	component := templates.Counters(server.GlobalStore.CounterImages)
	c.Header("Cache-Control", "no-cache")
	templ.Handler(component).ServeHTTP(c.Writer, c.Request)
}

func (web *webHandler) GETIndex(c *gin.Context) {
	server.GlobalStore.Lock()
	defer server.GlobalStore.Unlock()

	component := templates.Index()
	c.Header("Cache-Control", "no-cache")
	templ.Handler(component).ServeHTTP(c.Writer, c.Request)
}
