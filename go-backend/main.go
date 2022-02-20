package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hashicorp/go-hclog"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/moaabb/go-backend/models"
)

const version string = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	version string
}

type application struct {
	config  *config
	l       hclog.Logger
	DBModel *models.DBModel
}

type ApiStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

func main() {
	var cfg config
	l := hclog.Default()

	app := application{
		config: &cfg,
		l:      l,
	}

	flag.IntVar(&cfg.port, "port", 8080, "Port that the server will listen to")
	flag.StringVar(&cfg.env, "environment", "development", "development|production environment")
	flag.StringVar(&cfg.db.dsn, "dsn", "postgres://moab:example@localhost:8000/bookings?sslmode=disable", "Postgres Connection String")
	cfg.version = version

	db, err := OpenDB(cfg.db.dsn)
	if err != nil {
		l.Error("Couldn't connect to DB: ", err.Error())
		os.Exit(1)
	}
	l.Info("Connect to DB!")

	app.DBModel = models.NewDBModel(db)

	s := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Info("Server Listening on Port", cfg.port)
		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got Signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
