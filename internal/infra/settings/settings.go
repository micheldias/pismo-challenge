package settings

import "time"

var Settings settings

type settings struct {
	AppName               string        `default:"pismo-challenge-api-local"`
	TimeoutEndpoints      time.Duration `default:"15s"`
	ServerShutdownTimeout time.Duration `default:"40s"`
	Port                  string        `default:"8080"`
	DatabaseConnection    Database      `envconfig:"DATABASE"`
}

type Database struct {
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
	User     string `default:"meu_usuario"`
	Password string `default:"minha_senha_segura"`
	Dbname   string `default:"meu_banco"`
}
