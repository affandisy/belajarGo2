package notification

type Repository interface {
	SendEmail(toName, toEmail, subject, message string) (err error)
}
