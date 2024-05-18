package service

import (
	"projectOrder/internal/dto"
	"projectOrder/internal/pkg/cache"
	"projectOrder/internal/repository"
)

type Order interface {
	GetOrderById(orderID string) (*dto.Order, error)
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository, cache cache.Cache, orderExpTime int) *Service {
	return &Service{
		Order: NewOrderService(repos.Order, cache, orderExpTime),
	}
}
