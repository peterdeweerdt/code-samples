package core

type Mailer interface {
	SendEmail(template, toString, subject string, data interface{}) error
}
