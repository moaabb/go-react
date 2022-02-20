package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/moaabb/go-backend/data"
	"github.com/moaabb/go-backend/models"
)

func (app *application) StatusHandler(rw http.ResponseWriter, r *http.Request) {
	status := ApiStatus{
		Status:  "Active",
		Message: "API is Working",
		Version: "1.0.0",
	}

	err := data.ToJSON(rw, status, http.StatusOK, "api_status")
	if err != nil {
		app.l.Error(err.Error())
		return
	}
}

func (app *application) GetMovieByID(rw http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		app.l.Error("error parsing id, not an int", err.Error())
		response := data.GenericError{
			Message: "error parsing id, not an int",
		}
		data.ToJSON(rw, response, http.StatusBadRequest, "error")
		return
	}

	movie, err := app.DBModel.GetByID(id)
	if err != nil {
		response := data.GenericError{
			Message: err.Error(),
		}
		data.ToJSON(rw, response, http.StatusNotFound, "error")
		return
	}

	err = data.ToJSON(rw, movie, http.StatusOK, "movie")
	if err != nil {
		app.l.Error(err.Error())
		return
	}
}

func (app *application) GetAllMovies(rw http.ResponseWriter, r *http.Request) {
	movies, err := app.DBModel.GetAll()
	if err != nil {
		response := data.GenericError{
			Message: err.Error(),
		}
		data.ToJSON(rw, response, http.StatusNotFound, "error")
		return
	}

	err = data.ToJSON(rw, movies, http.StatusOK, "movies")
	if err != nil {
		app.l.Error(err.Error())
		return
	}
}

func (app *application) GetAllGenres(rw http.ResponseWriter, r *http.Request) {
	genres, err := app.DBModel.GetAllGenres()
	if err != nil {
		response := data.GenericError{
			Message: err.Error(),
		}
		data.ToJSON(rw, response, http.StatusNotFound, "error")
		return
	}

	err = data.ToJSON(rw, genres, http.StatusOK, "genres")
	if err != nil {
		app.l.Error(err.Error())
		return
	}
}

func (app *application) GetMovieByGenre(rw http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		response := data.GenericError{
			Message: "Couldn't resolve ID",
		}
		data.ToJSON(rw, response, http.StatusBadRequest, "error")
		return
	}

	movies, err := app.DBModel.GetAll(id)
	if err != nil {
		response := data.GenericError{
			Message: fmt.Sprintf("Not Found: %s", err.Error()),
		}
		data.ToJSON(rw, response, http.StatusNotFound, "error")
		return
	}

	genre, err := app.DBModel.GetGenreByID(id)
	if err != nil {
		response := data.GenericError{
			Message: fmt.Sprintf("Not Found: %s", err.Error()),
		}
		data.ToJSON(rw, response, http.StatusNotFound, "error")
		return
	}

	if len(movies) == 0 {
		movies = []*models.Movie{}
	}

	type response struct {
		Movies []*models.Movie `json:"movies"`
		Genre  *models.Genre   `json:"genre"`
	}

	data.ToJSON(rw, &response{
		Movies: movies,
		Genre:  genre,
	}, http.StatusOK, "movies_by_genre")
}

func (app *application) GetGenreByID(rw http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil {
		app.l.Error("error parsing id, not an int", err.Error())
		response := data.GenericError{
			Message: "error parsing id, not an int",
		}
		data.ToJSON(rw, response, http.StatusBadRequest, "error")
		return
	}

	genre, err := app.DBModel.GetGenreByID(id)
	if err != nil {
		response := data.GenericError{
			Message: err.Error(),
		}
		data.ToJSON(rw, response, http.StatusNotFound, "error")
		return
	}

	err = data.ToJSON(rw, genre, http.StatusOK, "genre")
	if err != nil {
		app.l.Error(err.Error())
		return
	}
}

func (app *application) EditMovie(rw http.ResponseWriter, r *http.Request) {
	var response models.Movie

	err := data.FromJSON(&response, r.Body)
	if err != nil {
		data.ToJSON(rw, &data.GenericError{
			Message: err.Error(),
		}, http.StatusBadRequest, "error")
		app.l.Error(err.Error())
		return
	}

	fmt.Println(response)

	data.ToJSON(rw, response, http.StatusOK, "response")
}

func (app *application) NotFound(rw http.ResponseWriter, r *http.Request) {
	response := data.GenericError{
		Message: "Not Found, Invalid URL",
	}

	data.ToJSON(rw, response, http.StatusNotFound, "error")
}
