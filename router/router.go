package router

import (
	"net/http"

	"github.com/yaninyzwitty/scylla-go-app/controller"
	"github.com/yaninyzwitty/scylla-go-app/middleware"
)

func NewRouter(controller *controller.Controller) *http.ServeMux {
	// create a new serve mux
	router := http.NewServeMux()

	// define middlewares here
	LoggingMiddleware := middleware.LoggingMiddleware
	corsMiddleware := middleware.CorsMiddleware

	// create a middleware chain
	middlewareChain := middleware.ChainMiddlewares(
		LoggingMiddleware,
		corsMiddleware,
	)

	// Define routes and wrap them with the middleware stack
	router.HandleFunc("POST /songs", func(w http.ResponseWriter, r *http.Request) {
		middlewareChain(http.HandlerFunc(controller.CreateSong)).ServeHTTP(w, r)
	})
	router.HandleFunc("PUT /songs/{id}", func(w http.ResponseWriter, r *http.Request) {
		middlewareChain(http.HandlerFunc(controller.UpdateSong)).ServeHTTP(w, r)
	})
	router.HandleFunc("GET /songs", func(w http.ResponseWriter, r *http.Request) {
		middlewareChain(http.HandlerFunc(controller.GetAllSongs)).ServeHTTP(w, r)
	})
	router.HandleFunc("GET /songs/{id}", func(w http.ResponseWriter, r *http.Request) {
		middlewareChain(http.HandlerFunc(controller.GetSong)).ServeHTTP(w, r)
	})
	router.HandleFunc("DELETE /songs/{id}", func(w http.ResponseWriter, r *http.Request) {
		middlewareChain(http.HandlerFunc(controller.DeleteSong)).ServeHTTP(w, r)
	})

	return router

}
