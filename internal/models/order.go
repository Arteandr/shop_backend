package models

import "time"

type Order struct {
	Id         int         `json:"id,omitempty" db:"id"`
	Items      []OrderItem `json:"items" db:"items"`
	UserId     int         `json:"userId" db:"user_id"`
	DeliveryId int         `json:"deliveryId" db:"delivery_id"`
	StatusId   int         `json:"statusId" db:"status_id"`
	Comment    string      `json:"comment" db:"comment"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
}

type OrderItem struct {
	Id       int `json:"itemId" binding:"required"`
	ColorId  int `json:"colorId" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}

type ServiceOrder struct {
	Id        int                `json:"id" db:"id"`
	Status    string             `json:"status" db:"status"`
	UserId    int                `json:"user_id" db:"user_id"`
	Items     []ServiceOrderItem `json:"items"`
	Delivery  Delivery           `json:"delivery"`
	Comment   string             `json:"comment" db:"comment"`
	CreatedAt time.Time          `json:"createdAt" db:"created_at"`
}

type ServiceOrderItem struct {
	Id          int     `json:"id" db:"item_id"`
	Name        string  `json:"name" db:"item_name"`
	Sku         string  `json:"sku" db:"item_sku"`
	Price       float64 `json:"price" db:"item_price"`
	ColorId     int     `json:"colorId" db:"color_id"`
	Description string  `json:"description" db:"item_description"`
	Quantity    int     `json:"quantity" db:"quantity"`
}

type OrderStatus struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
