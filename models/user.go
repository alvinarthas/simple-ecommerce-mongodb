package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is model for the customer
type User struct {
	ID                primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
	UserName          string              `bson:"user_name" json:"user_name"`
	FullName          string              `bson:"full_name" json:"full_name"`
	Email             string              `bson:"email" json:"email"`
	Password          string              `bson:"password" json:"-"`
	SocialID          string              `bson:"social_id" json:"-"`
	Provider          string              `bson:"provider" json:"-"`
	Avatar            string              `bson:"avatar" json:"avatar"`
	Phone             string              `bson:"phone" json:"phone"`
	Sex               string              `bson:"sex" json:"sex"`
	Role              int                 `bson:"role" json:"role"`
	HaveStore         int                 `bson:"have_store" json:"have_store"`
	IsActivate        int                 `bson:"is_active" json:"is_active"`
	VerificationToken string              `bson:"verification_token" json:"verification_token"`
	CreatedDate       time.Time           `json:"createdDate"`
	LastUpdate        primitive.Timestamp `json:"lastUpdate"`
	Store             Store
}
