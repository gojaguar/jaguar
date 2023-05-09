package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Builder is in charge of building a Server. It applies the builder design pattern to allow developers
// to build web servers according to their needs.
type Builder struct {
	router         chi.Router
	controllers    []Controller
	middlewares    []func(handler http.Handler) http.Handler
	signals        []os.Signal
	signalsChannel chan os.Signal
	port           uint16
	readTimeout    time.Duration
	writeTimeout   time.Duration
	idleTimeout    time.Duration
}

// Build builds the web server, it should be called at the very end of the builder chain.
func (builder *Builder) Build() *Server {
	builder.buildMiddlewares()
	builder.buildRoutes()
	builder.buildSignals()

	return &Server{
		http: builder.buildHTTP(),
		sigs: builder.signalsChannel,
	}
}

// buildRoutes is an internal method that processes all the routes before the final build call.
func (builder *Builder) buildRoutes() {
	if len(builder.controllers) == 0 {
		builder.controllers = defaultRoutes()
	}
	for _, ctrl := range builder.controllers {
		builder.router.Mount(ctrl.Namespace(), ctrl)
	}
}

// buildMiddlewares is an internal method that processes all the middlewares before the final build call.
func (builder *Builder) buildMiddlewares() {
	if len(builder.middlewares) == 0 {
		builder.middlewares = defaultMiddlewares()
	}
	for _, mw := range builder.middlewares {
		builder.router.Use(mw)
	}
}

// buildSignals is an internal method that processes all the signals before the final build call.
func (builder *Builder) buildSignals() {
	if len(builder.signals) == 0 {
		builder.signals = append(builder.signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	}
	builder.signalsChannel = make(chan os.Signal, 1)
	signal.Notify(builder.signalsChannel, builder.signals...)
}

// Controller allows the developer to specify a route or a set of routes that fulfill the Controller interface.
// This method allows multiple calls.
func (builder *Builder) Controller(ctrl Controller) *Builder {
	builder.controllers = append(builder.controllers, ctrl)
	return builder
}

// Middleware allows the developer to specify a middleware that will be executed before every handler.
// This method allows multiple calls.
func (builder *Builder) Middleware(mw func(handler http.Handler) http.Handler) *Builder {
	builder.middlewares = append(builder.middlewares, mw)
	return builder
}

// Signal allows the developer to specify an OS signal that will shut down the server once it listens to incoming
// HTTP requests. It's used for gracefully shutting down the final web server.
// This method allows multiple calls.
func (builder *Builder) Signal(signal os.Signal) *Builder {
	builder.signals = append(builder.signals, signal)
	return builder
}

// Port allows the developer to specify the HTTP port where to listen for incoming requests.
// This method allows a single call.
func (builder *Builder) Port(port uint16) *Builder {
	builder.port = port
	return builder
}

// Timeout allows the developer to specify the timeouts for writing, reading and idling once the web server
// is running. If no positive greater than zero values are specified or if this method is not called, there will be no
// timeout.
func (builder *Builder) Timeout(write time.Duration, read time.Duration, idle time.Duration) *Builder {
	builder.writeTimeout = write
	builder.readTimeout = read
	builder.idleTimeout = idle
	return builder
}

// buildHTTP builds the HTTP server, it defaults certain values if not specified to:
//   - port: 3030
func (builder *Builder) buildHTTP() *http.Server {
	if builder.port == 0 {
		builder.port = 3030
	}
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", builder.port),
		Handler:      builder.router,
		ReadTimeout:  builder.readTimeout,
		WriteTimeout: builder.writeTimeout,
		IdleTimeout:  builder.idleTimeout,
	}
}

// Controller groups a set of routes in a certain namespace.
//
//	The 'users' namespace can contain the following routes:
//	- GET /users/1
//	- POST /users
//	- PUT /users/1
//	- DELETE /users/1
type Controller interface {
	http.Handler
	// Namespace determines the route prefix used to expose the current controller.
	// For the users service, the namespace is users, giving the following route structure:
	// <namespace>/:id
	Namespace() string
}

// Server is used to listen to incoming HTTP requests until an OS signal is triggered.
type Server struct {
	// http contains a handler to the HTTP Server listening to incoming HTTP requests.
	http *http.Server
	// sigs contains the OS signals used to Shutdown the Server.
	sigs chan os.Signal
}

// ListenAndServe listens for incoming HTTP requests until an error occurs.
func (s *Server) ListenAndServe() error {
	errs := make(chan error, 1)

	go func(srv *http.Server, errs chan<- error) {
		if err := srv.ListenAndServe(); err != nil {
			errs <- err
			close(errs)
		}
	}(s.http, errs)

	var err error
	select {
	case sig := <-s.sigs:
		err = fmt.Errorf("signal %s triggered", sig)
		s.Shutdown()
	case e := <-errs:
		err = fmt.Errorf("error occurred: %w", e)
		s.Shutdown()
	}
	return err
}

// Shutdown attempts to shut down the current server until a timeout of a minute occurs.
// Calling Shutdown allows the web server to process pending HTTP requests.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	if err := s.http.Shutdown(ctx); err != nil {
		log.Println("Failed to shutdown HTTP server:", err)
	}
	cancel()
}

// defaultRoutes returns a set of default routes.
func defaultRoutes() []Controller {
	return []Controller{
		namespacedHandlerFunc{
			namespace: "",
			HandlerFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte("OK")); err != nil {
					http.Error(w, "Failed to write response", http.StatusInternalServerError)
				}
			}),
		},
	}
}

// namespacedHandlerFunc is used to expose an http.HandlerFunc in a certain namespace.
type namespacedHandlerFunc struct {
	namespace string
	http.HandlerFunc
}

// Namespace returns the namespace of the current http.HandlerFunc
func (p namespacedHandlerFunc) Namespace() string {
	return p.namespace
}

// defaultMiddlewares returns the middlewares that should use by default when initializing an HTTP server if no middlewares
// were provided.
func defaultMiddlewares() []func(handler http.Handler) http.Handler {
	return []func(handler http.Handler) http.Handler{
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	}
}
