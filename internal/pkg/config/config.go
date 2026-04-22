package config

import "time"

type Config struct {
	Environment string   `mapstructure:"environment"` // dev, staging, prod
	Services    Services `mapstructure:"services"`
	Database    Database `mapstructure:"database"`
}

type Services struct {
	Book BookService `mapstructure:"book"`
}

type BookService struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Database struct {
	Path            string        `mapstructure:"path"`              // 数据库文件路径
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // 最大连接数
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大生命周期
}
