package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/InternPulse/famtrust-backend-auth/internal/jwtmod"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandlers struct {
	models interfaces.Models
	mailer interfaces.Mailer
}

// @Summary		Login to FamTrust
// @Description	Login to FamTrust
// @Tags			User-Authentication
// @ID				login
// @Accept			json
// @Produce		json
// @Failure		401	{object}	loginSampleResponseError401
// @Failure		500	{object}	loginSampleResponseError500
// @Success		200	{object}	loginSampleResponse200
// @Param			credentials	body	loginRequest	true	"User Credentials"
// @Router			/login [post]
func (uh *UserHandlers) Login(c *gin.Context) {
	var loginPayload loginRequest

	err := c.ShouldBindBodyWithJSON(&loginPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, loginResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Invalid Credentials",
		})
		return
	}

	// Send sample email
	email := interfaces.EmailMsg{
		Subject:  "Worked",
		From:     "biz@famtrust.biz",
		To:       os.Getenv("SAMPLE_TO_EMAIL"),
		BodyText: "It worked! worked!!",
	}
	err = uh.mailer.SendMail(&email)
	if err != nil {
		fmt.Printf("Error from Mailer: %v", err)
	}

	// validate the user against the database
	user, err := uh.models.Users().GetUserByEmail(loginPayload.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, loginResponse{
			StatusCode: http.StatusUnauthorized,
			Status:     "error",
			Message:    "Invalid Credentials",
		})
		return
	}

	valid, err := uh.models.Users().PasswordMatches(user.PasswordHash, loginPayload.Password)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, loginResponse{
			StatusCode: http.StatusUnauthorized,
			Status:     "error",
			Message:    "Invalid Credentials",
		})
		return
	}

	token, err := jwtmod.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error has occured",
		})
	}

	payload := loginResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "User logged in successfully",
		Token:      token,
	}

	c.JSON(http.StatusOK, payload)
}

// @Summary		Validate User Login Token
// @Description	Validate User Login Token
// @Tags			User-Authentication
// @ID				validate
// @Accept			json
// @Produce		json
// @Failure		401	{object}	loginSampleResponseError401
// @Failure		500	{object}	loginSampleResponseError500
// @Success		200	{object}	validateSampleResponse200
// @Security		BearerAuth
// @Router			/validate [get]
func (uh *UserHandlers) Validate(c *gin.Context) {
	UserID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured",
		})
		return
	}

	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured",
		})
		return
	}

	// retrieve the user from the database
	user, err := uh.models.Users().GetUserByID(UserID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured",
		})
		return
	}

	permissions := func() []string {
		var list []string
		for _, perm := range user.Role.Permissions {
			list = append(list, perm.ID)
		}
		return list
	}

	// user payload
	role := role{
		Id:          user.Role.ID,
		Permissions: permissions(),
	}
	userPayload := loginUserData{
		Id:         user.ID,
		Email:      user.Email,
		Has2FA:     user.Has2FA,
		IsVerified: user.IsVerified,
		IsFreezed:  user.IsFreezed,
		LastLogin:  user.LastLogin,
		Role:       role,
	}

	payload := loginResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "User session is valid",
		Token:      token.(string),
		Data: map[string]loginUserData{
			"user": userPayload,
		},
	}

	c.JSON(http.StatusOK, payload)
}

func (uh *UserHandlers) Signup(c *gin.Context) {
}
