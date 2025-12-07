package user

type (
	User struct {
		ID              string
		Email           string
		Password        string
		Fullname        string
		Role            string
		IsEmailVerified bool
	}
)
