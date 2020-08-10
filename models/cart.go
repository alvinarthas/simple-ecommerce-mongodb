package models

// Cart collection
type Cart struct {
	Product     string `bson:"product" json:"product"`
	Qty         int    `bson:"qty" json:"qty"`
	Description string `bson:"description" json:"description"`
}
