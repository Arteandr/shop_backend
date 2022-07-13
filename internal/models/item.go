package models

type Item struct {
	Id          int     `json:"id,omitempty" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	CategoryId  int     `json:"categoryId" db:"category_id"`
	Tags        []Tag   `json:"tags,omitempty"`
	Colors      []Color `json:"colors,omitempty"`
	Sku         string  `json:"sku" db:"sku"`
}
