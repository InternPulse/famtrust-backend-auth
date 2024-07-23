package mailer

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
)

type Mailer struct{}

func NewMailer() interfaces.Mailer {
	return &Mailer{}
}
