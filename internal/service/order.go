package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"projectOrder/internal/dto"
	"projectOrder/internal/pkg/broker"
	"projectOrder/internal/pkg/cache"
	"projectOrder/internal/repository"
)

type orderService struct {
	repo         repository.Order
	cache        cache.Cache
	orderExpTime int
	broker       broker.IBroker
}

func NewOrderService(repo repository.Order, cache cache.Cache, orderExpTime int, broker broker.IBroker) Order {
	return &orderService{repo: repo, cache: cache, orderExpTime: orderExpTime, broker: broker}
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
		log.Println("Error retrieving order info from database:", err.Error())
		return nil, err
	}

	return infoOrder, err
}

func (o *orderService) StartConsumer(ctx context.Context, subject string) {
	orderChan := make(chan *dto.Order)

	go o.ProcessOrders(ctx, orderChan)

	subscription, err := o.broker.Subscribe(subject, func(msg *nats.Msg) {
		var infoOrder dto.Order
		fmt.Println(string(msg.Data))
		if err := json.Unmarshal(msg.Data, &infoOrder); err != nil {
			log.Println("Error decoding order info from broker:", err.Error())
			return
		}

		orderChan <- &infoOrder
	})
	if err != nil {
		log.Fatal("Error subscribing to order.info:", err.Error())
	}
	defer subscription.Unsubscribe()

	<-ctx.Done()
	close(orderChan)
}

func (o *orderService) ProcessOrders(ctx context.Context, orderChan <-chan *dto.Order) {
	for {
		select {
		case order := <-orderChan:
			if err := o.repo.SaveOrder(ctx, order); err != nil {
				log.Println("Error saving order to database:", err.Error())
			} else {
				log.Println("Order saved to database:", order)
			}
		case <-ctx.Done():
			log.Println("Context cancelled, stopping order processing")
			return
		}
	}
}
