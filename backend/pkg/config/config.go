package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig
	Session SessionConfig
	Log     LogConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type SessionConfig struct {
	Timeout        time.Duration
	CleanupInterval time.Duration
}

type LogConfig struct {
	Level  string
	Format string
}

var AppConfig Config

func InitConfig() {
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("session.timeout", 3600)
	viper.SetDefault("session.cleanup_interval", 300)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	viper.AutomaticEnv()

	AppConfig = Config{
		Server: ServerConfig{
			Host: viper.GetString("server.host"),
			Port: viper.GetInt("server.port"),
		},
		Session: SessionConfig{
			Timeout:        time.Duration(viper.GetInt("session.timeout")) * time.Second,
			CleanupInterval: time.Duration(viper.GetInt("session.cleanup_interval")) * time.Second,
		},
		Log: LogConfig{
			Level:  viper.GetString("log.level"),
			Format: viper.GetString("log.format"),
		},
	}
}
