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

type smtp struct {
	Host        string `env:"SMTP_HOST" default:"localhost"`
	Port        int    `env:"SMTP_PORT" default:"587"`
	Username    string `env:"SMTP_USERNAME" default:""`
	Password    string `env:"SMTP_PASSWORD" default:""`
	From        string `env:"SMTP_FROM" default:"noreply@pagoda.local"`
	FrontendURL string `env:"FRONTEND_URL" default:"http://localhost:5173"`
}
