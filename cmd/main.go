package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"projectOrder/internal/config"
	"projectOrder/internal/handler"
	"projectOrder/internal/pkg/broker"
	"projectOrder/internal/pkg/cache"
	"projectOrder/internal/pkg/psql"
	"projectOrder/internal/repository"
	"projectOrder/internal/server"
	"projectOrder/internal/service"
	"syscall"
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

	repos := repository.NewRepository(db)

	natsBroker, err := broker.NewBroker(cfg.BrokerConfig)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	cch := cache.NewCache(cfg.CacheConfig)
	domain := service.NewService(repos, cch, cfg.CacheConfig.OrderExpirationTime, natsBroker)
	routes := handler.NewHandler(domain)
	srv := server.NewServer(cfg.ServerConfig)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go domain.Order.StartConsumer(ctx, cfg.BrokerConfig.Subject)

	if err = srv.Run(routes.InitRoutes()); err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutdown signal received")
}
