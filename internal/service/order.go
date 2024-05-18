package service

import (
	"encoding/json"
	"log"
	"projectOrder/internal/dto"
	"projectOrder/internal/pkg/cache"
	"projectOrder/internal/repository"
)

type orderService struct {
	repo         repository.Order
	cache        cache.Cache
	orderExpTime int
}

func NewOrderService(repo repository.Order, cache cache.Cache, orderExpTime int) Order {
	return &orderService{repo: repo, cache: cache, orderExpTime: orderExpTime}
}

func (o *orderService) GetOrderById(orderID string) (*dto.Order, error) {
	infoOrder := new(dto.Order)

	infoOrderCache, err := o.cache.Get(orderID)
	if err == nil {
		if err = json.Unmarshal(infoOrderCache.Value, infoOrder); err != nil {
			log.Println(err.Error())
			return nil, err
		}

		return infoOrder, nil
	}

	infoOrder, err = o.repo.GetOrderById(orderID)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if err = o.saveOrderToCache(infoOrder); err != nil {
		return nil, err
	}

	return infoOrder, err
}

func (o *orderService) saveOrderToCache(infoOrder *dto.Order) error {
	cacheValue, err := json.Marshal(infoOrder)

	if err != nil {
		return err
	}

	return o.cache.Set(&cache.CacheItem{
		Key:     infoOrder.OrderUID,
		Value:   cacheValue,
		ExpTime: int32(o.orderExpTime),
	})
}
