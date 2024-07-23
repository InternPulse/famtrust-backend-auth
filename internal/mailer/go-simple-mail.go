package mailer

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	mail "github.com/xhit/go-simple-mail/v2"
)

func (m *Mailer) SendMail(email *interfaces.EmailMsg) error {
	newEmail := mail.NewMSG()
	newEmail.SetFrom(email.From).
		AddTo(email.To).
		SetSubject(email.Subject)

	newEmail.SetBody(mail.TextPlain, email.BodyText)

	server := newSMTPServer()
	client, err := server.Connect()
	if err != nil {
		return err
	}

	if err := newEmail.Send(client); err != nil {
		return err
	}

	return nil
}

func newSMTPServer() *mail.SMTPServer {
	mailer := mail.NewSMTPClient()

	// Get smtp details from env
	host := os.Getenv("SMTP_SERVER")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		log.Panicf("Env parse error: %v", err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	mailer.Host = host
	mailer.Port = int(port)
	mailer.Username = username
	mailer.Password = password

	mailer.Encryption = mail.EncryptionSTARTTLS
	mailer.ConnectTimeout = 30 * time.Second
	mailer.SendTimeout = 30 * time.Second

	return mailer
}
