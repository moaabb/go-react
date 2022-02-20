package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func NewDBModel(db *sql.DB) *DBModel {
	return &DBModel{
		DB: db,
	}
}

func (m *DBModel) GetByID(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var movie Movie

	query := `
		SELECT
			id, title, description, year, release_date, rating, runtime, mpaa_rating, created_at, updated_at
		FROM movies WHERE id = $1
	`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Rating,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	genres, err := m.GetMovieGenres(movie.ID)
	if err != nil {
		return nil, err
	}

	movie.MovieGenre = genres

	return &movie, nil
}

func (m *DBModel) GetAll(genres ...int) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var movies []*Movie

	where := ""
	if len(genres) > 0 {
		where = fmt.Sprintf(`
			where id in (select movie_id from movies_genres where genre_id = %d)
		`, genres[0])
	}

	query := `
		SELECT
			id, title, description, year, release_date, rating, runtime, mpaa_rating, created_at, updated_at
		FROM movies
	` + where

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var movie Movie
		err = rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Rating,
			&movie.Runtime,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		genres, err := m.GetMovieGenres(movie.ID)
		if err != nil {
			return nil, err
		}

		movie.MovieGenre = genres

		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (m *DBModel) GetGenreByID(id int) (*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT
			id, genre_name, created_at, updated_at
		FROM genres
		WHERE id = $1
	`

	var genre Genre

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&genre.ID,
		&genre.GenreName,
		&genre.CreatedAt,
		&genre.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &genre, nil

}

func (m *DBModel) GetAllGenres() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT
			id, genre_name, created_at, updated_at
		FROM genres ORDER BY genre_name
	`

	var genres []*Genre

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre Genre
		err = rows.Scan(
			&genre.ID,
			&genre.GenreName,
			&genre.CreatedAt,
			&genre.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, &genre)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}

func (m *DBModel) GetMovieGenres(id int) (map[int]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
			SELECT
				mg.id, mg.movie_id, mg.genre_id, g.genre_name, g.id
			FROM movies_genres mg
			LEFT JOIN genres g
			ON (mg.genre_id = g.id)
			WHERE mg.movie_id = $1
		`

	genres := make(map[int]string)
	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mg MovieGenre

		err = rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.Genre.GenreName,
			&mg.Genre.ID,
		)
		if err != nil {
			return nil, err
		}

		genres[mg.ID] = mg.Genre.GenreName

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}
