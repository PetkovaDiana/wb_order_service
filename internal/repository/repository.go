package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"projectOrder/internal/dto"
)

type Order interface {
	GetOrderById(orderID string) (*dto.Order, error)
	SaveOrder(ctx context.Context, order *dto.Order) error
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order: NewOrderRepo(db),
	}
}
