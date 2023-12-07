package app

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Order struct {
	Order_uid  string    `gorm:"index:idx_order_id,unique"`
	Order_body OrderBody `gorm:"type:jsonb"`
}

type OrderBody struct {
	Order_uid    string `json:"order_uid" binding:"required"`
	Track_number string `json:"track_number" binding:"required"`
	Entry        string `json:"entry" binding:"required"`
	Delivery     struct {
		Name    string `json:"name" binding:"required"`
		Phone   string `json:"phone" binding:"required"`
		Zip     string `json:"zip" binding:"required"`
		City    string `json:"city" binding:"required"`
		Address string `json:"address" binding:"required"`
		Region  string `json:"region" binding:"required"`
		Email   string `json:"email" binding:"required"`
	} `json:"delivery" binding:"required"`
	Payment struct {
		Transaction   string `json:"transaction" binding:"required"`
		Request_id    string `json:"request_id"`
		Currency      string `json:"currency" binding:"required"`
		Provider      string `json:"provider" binding:"required"`
		Amount        int    `json:"amount" binding:"required"`
		Payment_dt    int    `json:"payment_dt" binding:"required"`
		Bank          string `json:"bank" binding:"required"`
		Delivery_cost int    `json:"delivery_cost" binding:"required"`
		Goods_total   int    `json:"goods_total" binding:"required"`
		Custom_fee    int    `json:"custom_fee"`
	} `json:"payment" binding:"required"`
	Items []struct {
		Chrt_id      int    `json:"chrt_id" binding:"required"`
		Track_number string `json:"track_number" binding:"required"`
		Price        int    `json:"price" binding:"required"`
		Rid          string `json:"rid" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Sale         int    `json:"sale" binding:"required"`
		Size         string `json:"size" binding:"required"`
		Total_price  int    `json:"total_price" binding:"required"`
		Nm_id        int    `json:"nm_id" binding:"required"`
		Brand        string `json:"brand" binding:"required"`
		Status       int    `json:"status" binding:"required"`
	} `json:"items" binding:"required"`
	Locale             string `json:"locale" binding:"required"`
	Internal_signature string `json:"internal_signature"`
	Customer_id        string `json:"customer_id" binding:"required"`
	Delivery_service   string `json:"delivery_service" binding:"required"`
	Shardkey           string `json:"shardkey" binding:"required"`
	Sm_id              int    `json:"sm_id" binding:"required"`
	Date_created       string `json:"date_created" binding:"required"`
	Oof_shard          string `json:"oof_shard" binding:"required"`
}

/*func (order *Order) Scan(value interface{}) error {
	b, err := value.([]byte)
	if !err {
		log.Fatalf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &order)
}
func (order *Order) Value() (driver.Value, error) {
	return json.Marshal(order)
}*/

func (o OrderBody) Value() (driver.Value, error) {
	bytes, err := json.Marshal(o)
	return string(bytes), err
}

func (o *OrderBody) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []bytes")
	}
	return json.Unmarshal(bytes, o)
}
