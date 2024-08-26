package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ZeroBl21/go-further/internal/data"
	"github.com/ZeroBl21/go-further/internal/jsonlog"
	"github.com/ZeroBl21/go-further/internal/mailer"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

var buildTime string

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenCoons int
		maxIdleCoons int
		maxIdleTime  string
	}
	limiting struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	// Server
	flag.IntVar(&cfg.port, "port", 5173, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|stating|production)")

	// DB
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenCoons, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleCoons, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(
		&cfg.db.maxIdleTime,
		"db-max-idle-time",
		"15m",
		"PostgreSQL max connection idle time",
	)

	// Limiting
	flag.Float64Var(&cfg.limiting.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiting.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiting.enabled, "limiter-enabled", true, "Enable rate limiter")

	// SMTP
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")

	flag.StringVar(
		&cfg.smtp.username,
		"smtp-username",
		os.Getenv("GREENLIGHT_SMTP_USERNAME"),
		"SMTP username",
	)
	flag.StringVar(
		&cfg.smtp.password,
		"smtp-password",
		os.Getenv("GREENLIGHT_SMTP_PASSWORD"),
		"SMTP password",
	)
	flag.StringVar(
		&cfg.smtp.sender,
		"smtp-sender",
		"Greenlight <no-reply@greenlight.alexedwards.net>",
		"SMTP sender",
	)

	flag.Func(
		"cors-trusted-origins",
		"Trusted CORS origins (space separated)",
		func(val string) error {
			cfg.cors.trustedOrigins = strings.Fields(val)
			return nil
		},
	)

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	logger := jsonlog.New(os.Stdout, slog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(
			cfg.smtp.host,
			cfg.smtp.port,
			cfg.smtp.username,
			cfg.smtp.password,
			cfg.smtp.sender,
		),
	}

	if err := app.serve(); err != nil {
		logger.Fatal(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenCoons)
	db.SetMaxIdleConns(cfg.db.maxIdleCoons)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
