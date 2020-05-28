package collections

// User is model for the customer
type User struct {
	UserName          string `bson:"user_name" json:"user_name"`
	FullName          string `bson:"full_name" json:"full_name"`
	Email             string `bson:"email" json:"email"`
	Password          string `bson:"password" json:"password"`
	SocialID          string `bson:"social_id" json:"social_id"`
	Provider          string `bson:"provider" json:"provider"`
	Avatar            string `bson:"avatar" json:"avatar"`
	Role              int    `bson:"role" json:"role"`
	HaveStore         int    `bson:"have_store" json:"have_store"`
	IsActivate        int    `bson:"is_active" json:"is_active"`
	VerificationToken string `bson:"verification_token" json:"verification_token"`
	Store             *Store `bson:"store" json:",omitempty"`
}
