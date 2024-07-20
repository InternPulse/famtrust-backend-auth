package handlers

import (
	"net/http"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/InternPulse/famtrust-backend-auth/internal/jwtmod"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandlers struct {
	models interfaces.Models
}

func (uh *UserHandlers) Login(c *gin.Context) {
	var loginPayload loginRequest

	err := c.ShouldBindBodyWithJSON(&loginPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, loginResponse{
			StatusCode:    http.StatusBadRequest,
			StatusMessage: "error",
		})
		return
	}

	// validate the user against the database
	user, err := uh.models.Users().GetUserByEmail(loginPayload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "error",
		})
		return
	}

	valid, err := uh.models.Users().PasswordMatches(user, loginPayload.Password)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, loginResponse{
			StatusCode:    http.StatusUnauthorized,
			StatusMessage: "error",
		})
		return
	}

	token, err := jwtmod.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "error",
		})
	}

	payload := loginResponse{
		StatusCode:    http.StatusOK,
		StatusMessage: "success",
		Token:         token,
	}

	c.JSON(http.StatusOK, payload)
}

func (uh *UserHandlers) Validate(c *gin.Context) {
	UserID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "error",
		})
		return
	}

	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "error",
		})
		return
	}

	// retrieve the user from the database
	user, err := uh.models.Users().GetUserByID(UserID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode:    http.StatusInternalServerError,
			StatusMessage: "error",
		})
		return
	}

	// user payload
	Role := Role{
		ID:          user.Role.ID,
		Permissions: user.Role.Permissions,
	}
	userPayload := loginUserData{
		ID:    user.ID,
		Email: user.Email,
		Role:  Role,
	}

	payload := loginResponse{
		StatusCode:    http.StatusOK,
		StatusMessage: "success",
		Token:         token.(string),
		Data: map[string]loginUserData{
			"user": userPayload,
		},
	}

	c.JSON(http.StatusOK, payload)
}

func (uh *UserHandlers) Signup(c *gin.Context) {
}
