package models

type Item struct {
	Id          int      `json:"id,omitempty" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	Category    Category `json:"category"`
	Images      []Image  `json:"images,omitempty"`
	Tags        []Tag    `json:"tags,omitempty"`
	Colors      []Color  `json:"colors,omitempty"`
	Price       float64  `json:"price" db:"price"`
	Sku         string   `json:"sku" db:"sku"`
}
