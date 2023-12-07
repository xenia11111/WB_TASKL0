package repository

import (
	app "github.com/xenia11111/WB_TASKL0"
	"gorm.io/gorm"
)

type OrderCRUD interface {
	GetById(id string) (*app.OrderBody, error)
	Create(order app.Order) error
	GetAll() (*[]app.Order, error)
	GetBulk(countOfRecords int) (*[]app.Order, error)
}

type Repository struct {
	OrderCRUD
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		OrderCRUD: NewToDoOrderPostgres(db),
	}
}
