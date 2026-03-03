package config

import "time"

type server struct {
	Host        string        `env:"HOST" default:"0.0.0.0"`
	Port        int           `env:"PORT" default:"3000"`
	Secret      string        `env:"SECRET" default:"pagoda-secret"`
	Debug       bool          `env:"DEBUG" default:"false"`
	CorsOrigins string        `env:"CORS_ORIGINS" default:"*"`
	TokenExpiry time.Duration `env:"TOKEN_EXPIRY" default:"720h"`
}

type database struct {
	Driver string `env:"DB_DRIVER" default:"sqlite"`
	DSN    string `env:"DSN" default:"pagoda.db"`
}
