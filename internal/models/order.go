package models

type Order struct {
	Id         int         `json:"id,omitempty" db:"id"`
	Items      []OrderItem `json:"items" db:"items"`
	UserId     int         `json:"userId" db:"user_id"`
	DeliveryId int         `json:"deliveryId" db:"delivery_id"`
	StatusId   int         `json:"statusId" db:"status_id"`
}

type OrderItem struct {
	Id       int `json:"itemId" binding:"required"`
	ColorId  int `json:"colorId" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}
