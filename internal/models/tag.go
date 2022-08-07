package models

type Tag struct {
	Id     int    `json:"id,omitempty" db:"id"`
	ItemId int    `json:"itemId" db:"item_id"`
	Name   string `json:"name" db:"name"`
}
