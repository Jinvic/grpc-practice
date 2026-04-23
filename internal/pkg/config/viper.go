package config

import (
	"fmt"

	"github.com/samber/do/v2"
	"github.com/spf13/viper"
)

func GetConfig(i do.Injector) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	setDefaultConfig(viper.GetViper())

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

func setDefaultConfig(v *viper.Viper) {
	v.SetDefault("database.path", "./data/bookstore.db")
	v.SetDefault("database.max_open_conns", 1)
	v.SetDefault("database.max_idle_conns", 1)
	v.SetDefault("database.conn_max_lifetime", "5m")

	v.SetDefault("logging.level", "debug")
	v.SetDefault("logging.format", "text")
	v.SetDefault("logging.output", "stdout")
	v.SetDefault("logging.add_source", false)
	v.SetDefault("logging.max_size", 100)
	v.SetDefault("logging.max_age", 7)
	v.SetDefault("logging.max_backups", 10)
	v.SetDefault("logging.compress", false)
	v.SetDefault("logging.local_time", false)

	v.SetDefault("services.book.host", "localhost")
	v.SetDefault("services.book.port", 8081)
	v.SetDefault("services.book.log_file", "./logs/book_service.log")
}
