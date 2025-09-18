package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type application struct {
	logger *log.Logger
	db     *sql.DB
	rdb    *redis.Client
}

func main() {
	logger := log.New(os.Stdout, "[blog] ", log.LstdFlags)

	// Read config from env vars (set in docker-compose)
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSL := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, dbSSL)

	DB, err := openDB(dsn)
	if err != nil {
		logger.Fatal("Failed to connect to Database: ", err)
	}

	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		logger.Fatal("Can't create migration connection", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrations", "postgres", driver)
	if err != nil {
		logger.Fatal("Can't run migration", err)
	}

	m.Up() // or m.Steps(2) if you want to explicitly set the number of migrations to run

	// Redis connection
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	rdb, err := openRedis(redisAddr)
	if err != nil {
		logger.Fatal("Failed to connect to Redis: ", err)
	}

	me := application{
		logger: logger,
		db:     DB,
		rdb:    rdb,
	}

	logger.Println("Starting API server on :8000")
	if err := me.serve(); err != nil {
		logger.Fatal(err)
	}
}

func (me *application) serve() error {
	srv := &http.Server{
		Addr:         ":8000",
		Handler:      me.routes(),
		ErrorLog:     me.logger,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return srv.ListenAndServe()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func openRedis(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
