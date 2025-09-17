package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type application struct {
	logger *log.Logger
	db     *sql.DB
}

func main() {

	logger := log.New(os.Stdout, "[blog]", log.Flags())
	dsn := "postgresql://postgres:password@localhost:5432/blog_db?sslmode=disable"
	DB, err := openDB(dsn)
	if err != nil {
		logger.Fatal("Failed to connect to Database", err)
	}

	me := application{
		logger: logger,
		db:     DB,
	}

	me.serve()
}

func (me *application) serve() error {
	// Declare an HTTP server using the same settings as in our main() function.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8000),
		Handler:      me.routes(),
		ErrorLog:     me.logger,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return srv.ListenAndServe()
}

func openDB(dsn string) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config struct.
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// The values is made-up, profile your postgres in prod
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database,
	// passing in the context we created above as a parameter.
	// If connection couldn't be established successfully within the 5-second deadline,
	// then this will return an error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}
