package config

import (
	"github.com/BurntSushi/toml"
)

type Postgres struct {
	Host     string
	Port     int64
	User     string
	Password string
	Dbname   string
}

type MetaData struct {
	Description string
	Severity    string
}

type TokenScanner struct {
	Token    string
	Type     string
	RuleId   string
	MetaData MetaData
}

type Config struct {
	Postgres     Postgres
	TokenScanner []TokenScanner
}

func ReadConfig() (*Config, error) {
	var config Config
	_, err := toml.DecodeFile("./config.toml", &config)

	return &config, err
}
