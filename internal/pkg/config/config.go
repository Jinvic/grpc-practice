package config

import "time"

type Config struct {
	Environment string   `mapstructure:"environment"` // dev, staging, prod
	Database    Database `mapstructure:"database"`
	Logging     Logging  `mapstructure:"logging"`
	Services    Services `mapstructure:"services"`
}

type Services struct {
	Book BookService `mapstructure:"book"`
}

type BookService struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	LogFile string `mapstructure:"log_file"`
}

type Database struct {
	File            string        `mapstructure:"file"`              // 数据库文件路径
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // 最大连接数
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大生命周期
}

type Logging struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Format     string `mapstructure:"format"`      // json, text
	Output     string `mapstructure:"output"`      // stdout, file
	AddSource  bool   `mapstructure:"add_source"`  // 是否添加来源信息
	MaxSize    int    `mapstructure:"max_size"`    // 日志文件最大大小
	MaxAge     int    `mapstructure:"max_age"`     // 日志文件最大保存时间
	MaxBackups int    `mapstructure:"max_backups"` // 日志文件最大备份数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩日志文件
	LocalTime  bool   `mapstructure:"local_time"`  // 是否使用本地时间
}
