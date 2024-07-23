package interfaces

type Mailer interface {
	SendMail(email *EmailMsg) error
}

type EmailMsg struct {
	Subject  string
	From     string
	To       string
	BodyText string
}
