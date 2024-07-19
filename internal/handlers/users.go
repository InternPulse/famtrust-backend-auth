package handlers

import (
	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	models interfaces.Models
}

func (uh *UserHandlers) Signup(c *gin.Context) {
}

func (uh *UserHandlers) Login(c *gin.Context) {
}

func (uh *UserHandlers) Validate(c *gin.Context) {
}
