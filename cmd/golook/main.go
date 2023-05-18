package main

import (
	"context"
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mokan-r/golook/internal/config"
	"github.com/mokan-r/golook/internal/executor"
	"github.com/mokan-r/golook/internal/monitor"
	"github.com/mokan-r/golook/pkg/db"
	"github.com/mokan-r/golook/pkg/db/postgresql"
	"log"
	"os"
)

var (
	configPath = flag.String("config", "config.yml", "path to the configuration file")
	dsn        = flag.String("dsn", "postgres://golook:password@localhost:5432/golook", "Postgres data source name")
)

func init() {
	flag.Parse()
}

type application struct {
	db       db.DB
	monitor  monitor.Monitor
	executor executor.Executor
	logError *log.Logger
	logInfo  *log.Logger
}

func main() {
	logError := log.New(
		os.Stdout,
		"ERROR\t",
		log.Ldate|log.Ltime,
	)
	logInfo := log.New(
		os.Stdout,
		"INFO\t",
		log.Ldate|log.Ltime,
	)

	cfg, err := config.ReadConfig(*configPath, &config.DefaultReader{})
	if err != nil {
		logError.Fatal("Error reading config")
	}
	Database, err := openDB(*dsn)
	if err != nil {
		logError.Fatal(err)
	}
	defer Database.Close()
	eventChan := make(chan monitor.EventTrigger)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logError.Fatal(err)
	}
	mon, err := monitor.New(watcher, cfg, eventChan)
	if err != nil {
		logError.Fatal(err)
	}

	app := application{
		db:       &postgresql.PostgreSQL{DB: Database},
		monitor:  mon,
		executor: executor.New(eventChan, cfg),
		logError: logError,
		logInfo:  logInfo,
	}

	app.monitor.Start()
	app.executor.Start()

	for {
		select {
		case command, ok := <-app.executor.Commands():
			if !ok {
				app.logError.Println("Something went wrong while getting commands from channel")
			}
			err := app.db.Insert(command)
			if err != nil {
				app.logError.Println("Something went wrong while inserting to database:", err)
			}
		}
	}
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return pool, nil
}
