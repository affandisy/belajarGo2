package user

type Repository interface {
	Create(user User) (err error)
	GetByEmail(email string) (user User, err error)
	UpdateEmailVerification(user User) (err error)
}
