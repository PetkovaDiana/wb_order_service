package main

import (
	"projectOrder/internal/config"
	"projectOrder/internal/handler"
	"projectOrder/internal/pkg/cache"
	"projectOrder/internal/pkg/psql"
	"projectOrder/internal/repository"
	"projectOrder/internal/server"
	"projectOrder/internal/service"
)

func main() {
	cfg, err := config.NewAppConfig()

	if err != nil {
		panic(err)
	}

	db, err := psql.NewDB(cfg.DBConfig)

	if err != nil {
		panic(err)
	}

	cch := cache.NewCache(cfg.CacheConfig)
	repos := repository.NewRepository(db)
	domain := service.NewService(repos, cch, cfg.CacheConfig.OrderExpirationTime)
	routes := handler.NewHandler(domain)
	srv := server.NewServer(cfg.ServerConfig)

	if err = srv.Run(routes.InitRoutes()); err != nil {
		panic(err)
	}
}
