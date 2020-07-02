package models

import "time"

// Payment collection
type Payment struct {
	PaymentID string    `bson:"payment_id" json:"payment_id"`
	Method    string    `bson:"method" json:"method"`
	Amount    int64     `bson:"amount" json:"amount"`
	Date      time.Time `bson:"date" json:"date"`
	File      string    `bson:"file" json:"file"`
}
