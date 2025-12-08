package user

type (
	User struct {
		ID              string `bson:"user_id"`
		Email           string
		Password        string
		Fullname        string
		Role            string
		IsEmailVerified bool `bson:"is_email_verified"`
	}
)
