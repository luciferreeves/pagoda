package config

type server struct {
	Host   string `env:"HOST" default:"0.0.0.0"`
	Port   int    `env:"PORT" default:"3000"`
	Secret string `env:"SECRET" default:"pagoda-secret"`
	Debug  bool   `env:"DEBUG" default:"false"`
}

type database struct {
	Driver string `env:"DB_DRIVER" default:"sqlite"`
	DSN    string `env:"DSN" default:"pagoda.db"`
}
