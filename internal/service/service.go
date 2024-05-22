package service

import (
	"context"
	"projectOrder/internal/dto"
	"projectOrder/internal/pkg/broker"
	"projectOrder/internal/pkg/cache"
	"projectOrder/internal/repository"
)

type Order interface {
	GetOrderById(orderID string) (*dto.Order, error)
	StartConsumer(ctx context.Context, subject string)
	ProcessOrders(ctx context.Context, orderChan <-chan *dto.Order)
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository, cache cache.Cache, orderExpTime int, broker broker.IBroker) *Service {
	return &Service{
		Order: NewOrderService(repos.Order, cache, orderExpTime, broker),
	}
}
