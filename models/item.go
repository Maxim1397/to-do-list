package models

//Item model for items table
type Item struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}
