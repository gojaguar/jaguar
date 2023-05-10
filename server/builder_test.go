package server

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestServerBuilder_DefaultConfig(t *testing.T) {
	var builder Builder

	srv := builder.Build()
	assert.NotNil(t, srv.http)
	assert.NotNil(t, srv.sigs)
	assert.NotNil(t, srv.http.Handler)
	assert.Equal(t, ":3030", srv.http.Addr)
}

func TestServerBuilder_WithPort(t *testing.T) {
	var builder Builder

	srv := builder.Port(6060).Build()
	assert.NotNil(t, srv.http)
	assert.NotNil(t, srv.sigs)
	assert.Equal(t, ":6060", srv.http.Addr)
}

func TestServerBuilder_WithTimeout(t *testing.T) {
	var builder Builder

	const write = 30 * time.Second
	const read = 15 * time.Second
	const idle = 5 * time.Second

	srv := builder.Timeout(
		30*time.Second, // write
		15*time.Second, // read
		5*time.Second,  // idle
	).Build()

	assert.NotNil(t, srv.http)
	assert.NotNil(t, srv.sigs)
	assert.Equal(t, write, srv.http.WriteTimeout)
	assert.Equal(t, read, srv.http.ReadTimeout)
	assert.Equal(t, idle, srv.http.IdleTimeout)
}

func TestServerBuilder_WithMiddleware(t *testing.T) {
	var builder *Builder
	builder = new(Builder)

	builder = builder.Middleware(middleware.AllowContentType("application/json"))

	assert.NotEmpty(t, builder.middlewares)
	assert.Len(t, builder.middlewares, 1)

	// When using the default config, it should contain 4 middlewares.
	builder = new(Builder)
	_ = builder.Build()

	assert.NotEmpty(t, builder.middlewares)
	assert.Len(t, builder.middlewares, 4)
}

func TestServerBuilder_WithController(t *testing.T) {
	var builder *Builder
	builder = new(Builder)

	builder = builder.Controller(namespacedHandlerFunc{
		namespace:   "test",
		HandlerFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	})

	assert.NotEmpty(t, builder.controllers)
	assert.Len(t, builder.controllers, 1)
	assert.Equal(t, "test", builder.controllers[0].Namespace())

	// When using the default config, it should contain 1 controller.
	builder = new(Builder)
	_ = builder.Build()

	assert.NotEmpty(t, builder.controllers)
	assert.Len(t, builder.controllers, 1)
}

func TestServerBuilder_WithSignal(t *testing.T) {
	var builder *Builder
	builder = new(Builder)

	builder = builder.Signal(syscall.SIGINT)

	assert.NotEmpty(t, builder.signals)
	assert.Len(t, builder.signals, 1)
	assert.Equal(t, syscall.SIGINT, builder.signals[0])

	// When using the default config, it should contain 3 signals.
	builder = new(Builder)
	_ = builder.Build()

	assert.NotEmpty(t, builder.signals)
	assert.Len(t, builder.signals, 3)
}
