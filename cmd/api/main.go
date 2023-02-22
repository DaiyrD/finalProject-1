package main

import (
	"context"
	"database/sql"
	"finalProjectAdvancedP/internal/data"
	"finalProjectAdvancedP/internal/jsonlog"
	"finalProjectAdvancedP/internal/mailer"
	"flag"
	"os"
	"sync"
	"time"
)

const version = "1.0.0"
const DATABASE_URL = "postgres://postgres:7151@localhost/bookstore?sslmode=disable"

// application's configuration settings (port, environment, DB, smtp server settings)
type config struct {
	port int
	env  string
	db   struct {
		dsn string // database source name
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

// application struct
// needs to be done
type application struct {
	models data.Models
	config config
	mailer mailer.Mailer
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

// starting point of our application
func main() {
	// declaring config instance
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "API server port")
	flag.StringVar(&cfg.env, "environment", "development", "Environment (development)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", DATABASE_URL, "PostgreSQL dsn")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.office365.com", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "211140@astanait.edu.kz", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "Aitu2021!", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "211140@astanait.edu.kz", "SMTP sender")

	flag.Parse()

	//logger

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	// declare an instance of our application
	// Use the data.NewModels() function to initialize a Models struct, passing in the
	// connection pool as a parameter.
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		// Initialize a new Mailer instance using the settings from the command line
		// flags, and add it to the application struct.
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
