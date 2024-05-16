package config

import (
	"main.go/projectOrder/internal/pkg/memcached"
	"main.go/projectOrder/internal/pkg/psql"
)

type App struct {
	DBConfig        psql.Config
	MemcachedConfig memcached.Config
}
