package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

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
		err = e
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
