package service

import (
	app "github.com/xenia11111/WB_TASKL0"
	"github.com/xenia11111/WB_TASKL0/pkg/cache"
	"github.com/xenia11111/WB_TASKL0/pkg/repository"
)

type OrderCRUD interface {
	GetById(id string) (*app.OrderBody, error)
	Create(order app.Order) error
}

type Service struct {
	OrderCRUD
}

func NewService(repos *repository.Repository, ch *cache.Cache) *Service {
	return &Service{

		OrderCRUD: NewOrderCRUDService(repos.OrderCRUD, ch.OrderCRUD),
	}
}
