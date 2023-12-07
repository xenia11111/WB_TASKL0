package repository

import (
	app "github.com/xenia11111/WB_TASKL0"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewToDoOrderPostgres(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetById(id string) (*app.OrderBody, error) {

	order_body := app.OrderBody{}
	err := r.db.Raw("SELECT order_body FROM orders WHERE order_uid = ?", id).Row().Scan(&order_body)

	if err != nil {
		return nil, err
	}

	return &order_body, nil
}

func (r *OrderRepository) Create(order app.Order) error {

	result := r.db.Create(&order)

	return result.Error
}

func (r *OrderRepository) GetAll() (*[]app.Order, error) {

	var res []app.Order

	result := r.db.Find(&res)

	if result.Error != nil {
		return nil, result.Error
	}

	return &res, nil
}

func (r *OrderRepository) GetBulk(countOfRecords int) (*[]app.Order, error) {

	var res []app.Order

	result := r.db.Limit(int(countOfRecords)).Find(&res)

	if result.Error != nil {
		return nil, result.Error
	}

	return &res, nil
}
