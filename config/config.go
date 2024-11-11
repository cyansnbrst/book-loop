package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const Version = "1.0.0"

type Config struct {
	Port       int
	Env        string
	PostgreSQL struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Limiter struct {
		RPS     float64
		Burst   int
		Enabled bool
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("error loading .env file")
	}

	var cfg Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.PostgreSQL.DSN, "db-dsn", os.Getenv("BOOKLOOP_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.PostgreSQL.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.PostgreSQL.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.PostgreSQL.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.RPS, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.SMTP.Host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.SMTP.Port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.SMTP.Username, "smtp-username", "2ad8d51e642dc5", "SMTP username")
	flag.StringVar(&cfg.SMTP.Password, "smtp-password", "45ef22500fbb6a", "SMTP password")
	flag.StringVar(&cfg.SMTP.Sender, "smtp-sender", "BookLoop <no-reply@bookloop.net>", "SMTP sender")

	flag.Parse()

	return &cfg
}
