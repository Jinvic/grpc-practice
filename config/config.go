package config

type Config struct {
	Environment string   `mapstructure:"environment"` // dev, staging, prod
	Services    Services `mapstructure:"services"`
}

type Services struct {
	Book BookService `mapstructure:"book"`
}

type BookService struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
