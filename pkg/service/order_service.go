package service

import (
	app "github.com/xenia11111/WB_TASKL0"
	"github.com/xenia11111/WB_TASKL0/pkg/cache"
	"github.com/xenia11111/WB_TASKL0/pkg/repository"
)

type OrderService struct {
	repo repository.OrderCRUD
	ch   cache.OrderCRUD
}

func NewOrderCRUDService(repo repository.OrderCRUD, cache cache.OrderCRUD) *OrderService {
	return &OrderService{
		repo: repo,
		ch:   cache,
	}
}

func (s *OrderService) GetById(id string) (*app.OrderBody, error) {

	order, found := s.ch.Get(id)

	if found {
		return order, nil
	}

	order, err := s.repo.GetById(id)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) Create(order app.Order) error {

	err := s.repo.Create(order)
	if err != nil {
		return err
	}
	s.ch.Set(order.Order_uid, order, 0)
	return nil
}
