package collections

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
Store is for the sellers,
User can have more than one stores
*/

// Store model
type Store struct {
	ID                primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name              string             `bson:"name" json:"name"`
	UserName          string             `bson:"user_name" json:"user_name"`
	Adress            string             `bson:"address" json:"address"`
	Email             string             `bson:"email" json:"email"`
	Phone             string             `bson:"phone" json:"phone"`
	Avatar            string             `bson:"avatar" json:"avatar"`
	IsActivate        int                `bson:"is_active,omitempty" json:"is_active,omitempty"`
	VerificationToken string             `bson:"verification_token" json:"verification_token"`
	Products          []string           `bson:"products" json:"products"`
}