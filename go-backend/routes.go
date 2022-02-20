package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() chi.Router {
	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			app.l.Info(r.Method, r.RequestURI)

			next.ServeHTTP(rw, r)
		})
	})
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/status", app.StatusHandler)
	r.Get("/v1/movies", app.GetAllMovies)
	r.Get("/v1/movies/genres/{id}", app.GetMovieByGenre)
	r.Get("/v1/movies/{id}", app.GetMovieByID)
	r.Post("/v1/movies/editmovie", app.EditMovie)
	r.Get("/v1/genres", app.GetAllGenres)
	r.NotFound(app.NotFound)

	return r
}
