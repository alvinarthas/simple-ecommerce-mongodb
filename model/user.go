package model

// User is model for the customer
type User struct {
	UserName          string `bson:"user_name" json:"user_name"`
	FullName          string `bson:"full_name" json:"full_name"`
	Email             string `bson:"email" json:"email"`
	Password          string `bson:"password" json:"password"`
	SocialID          string `bson:"social_id" json:"social_id"`
	Provider          string `bson:"provider" json:"provider"`
	Avatar            string `bson:"avatar" json:"avatar"`
	Role              *bool  `bson:"role,omitempty" json:"role,omitempty"`
	HaveStore         *bool  `bson:"have_store,omitempty" json:"have_store,omitempty"`
	IsActivate        *bool  `bson:"is_active,omitempty" json:"is_active,omitempty"`
	VerificationToken string `bson:"verification_token" json:"verification_token"`
	// Store             Store  // to show that customer can have one store
}
