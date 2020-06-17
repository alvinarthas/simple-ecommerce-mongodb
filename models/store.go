package models

/*
User can be the merchant, only available to one store
*/

// Store model
type Store struct {
	Name              string    `bson:"name" json:"name"`
	UserName          string    `bson:"user_name" json:"user_name"`
	Adress            string    `bson:"address" json:"address"`
	Email             string    `bson:"email" json:"email"`
	Phone             string    `bson:"phone" json:"phone"`
	Avatar            string    `bson:"avatar" json:"avatar"`
	IsActivate        int       `bson:"is_active,omitempty" json:"is_active,omitempty"`
	VerificationToken string    `bson:"verification_token" json:"verification_token"`
	Products          [0]string `bson:"products" json:"products"`
}
