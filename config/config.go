package config

import (
	"github.com/BurntSushi/toml"
)

type PostgresConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	Dbname   string
}

type GitConfig struct {
	CloneLocation string
}

type QueueConfig struct {
	QueueSize int64
}

type MetaData struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type ScanData struct {
	Token  string `json:"-"`
	Type   string `json:"type"`
	RuleId string `json:"rule_id"`
}

type TokenScannerConfig struct {
	ScanData *ScanData
	MetaData *MetaData `json:"metadata"`
}

type Config struct {
	Postgres     PostgresConfig
	TokenScanner []TokenScannerConfig
	Queue        QueueConfig
	Git          GitConfig
}

func ReadConfig() (*Config, error) {
	var config Config
	_, err := toml.DecodeFile("./config.toml", &config)

	return &config, err
}
