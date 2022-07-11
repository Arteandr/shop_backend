package models

type Item struct {
	Id          int      `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	CategoryId  int      `json:"categoryId" db:"category_id"`
	Tags        []string `json:"tags" db:"tags"`
	Colors      []Color  `json:"colors,omitempty"`
	Sku         string   `json:"sku" db:"sku"`
}
