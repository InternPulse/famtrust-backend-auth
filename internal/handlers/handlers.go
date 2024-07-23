package handlers

import "github.com/InternPulse/famtrust-backend-auth/internal/interfaces"

type Handlers struct {
	users         interfaces.UserHandlers
	verifications interfaces.VerificationHandlers
}

func (h *Handlers) Users() interfaces.UserHandlers {
	return h.users
}

func (h *Handlers) Verifications() interfaces.VerificationHandlers {
	return h.verifications
}

func NewHandler(models interfaces.Models, mailer interfaces.Mailer) interfaces.Handlers {
	return &Handlers{
		users:         &UserHandlers{models: models, mailer: mailer},
		verifications: &VerificationHandlers{models: models, mailer: mailer},
	}
}
