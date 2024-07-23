package handlers

import "github.com/InternPulse/famtrust-backend-auth/internal/interfaces"

type Handlers struct {
	users interfaces.UserHandlers
}

func (h *Handlers) Users() interfaces.UserHandlers {
	return h.users
}

func NewHandler(models interfaces.Models, mailer interfaces.Mailer) interfaces.Handlers {
	return &Handlers{
		users: &UserHandlers{models: models, mailer: mailer},
	}
}
