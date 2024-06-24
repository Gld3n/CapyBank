package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type configuration struct {
	dsn  string
	port string
}

type application struct {
	logger       *slog.Logger
	transactions *TransactionModel
}

func main() {
	var cfg configuration

	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("postgres-dsn"), "the connection string for the database")
	flag.StringVar(&cfg.port, "port", ":8000", "the address to mount the server onto")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(&cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}(db)

	app := &application{
		logger:       logger,
		transactions: &TransactionModel{DB: db},
	}

	logger.Info("starting server", slog.String("port", cfg.port))

	if err := http.ListenAndServe(cfg.port, app.routes()); err != nil {
		logger.Error("error starting server")
	}
}

func openDB(dsn *string) (*sql.DB, error) {
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
