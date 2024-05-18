package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"projectOrder/internal/pkg/cache"
	"projectOrder/internal/pkg/psql"
	"projectOrder/internal/server"
)

const (
	envPath = ".env"
)

type App struct {
	DBConfig     *psql.Config   `yaml:"db"`
	CacheConfig  *cache.Config  `yaml:"cache"`
	ServerConfig *server.Config `yaml:"server"`
}

func NewAppConfig() (*App, error) {
	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	cfgApp := new(App)

	if err := cleanenv.ReadConfig(os.Getenv("CONFIG_PATH"), cfgApp); err != nil {
		return nil, err
	}

	return cfgApp, nil
}
