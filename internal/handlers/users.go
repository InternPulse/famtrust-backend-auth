package handlers

import (
	"fmt"
	"net/http"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/InternPulse/famtrust-backend-auth/internal/jwtmod"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandlers struct {
	models interfaces.Models
	mailer interfaces.Mailer
}

// @Summary		Login to FamTrust (Supports 2FA by Email)
// @Description	Login to FamTrust (Supports 2FA by Email)
// @Tags			User-Authentication
// @ID				login
// @Accept			json
// @Produce		json
// @Failure		401	{object}	loginSampleResponseError401
// @Failure		500	{object}	loginSampleResponseError500
// @Success		200	{object}	loginSampleResponse200
// @Param			Credentials	body	loginRequest	true	"User Credentials"
// @Param			2FACode	query	string	false	"User 2FA Code"
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

	if user.Has2FA {

		twoFACode := c.Query("2FACode")
		// If no 2FA code is passed in, generate one
		if twoFACode == "" {

			verCode := interfaces.VerCode{
				UserID: user.ID,
				Type:   "2fa",
			}
			// NOTE: Created tokens are invalid once another is created
			// Only the latest/last created token is used
			err := uh.models.VerCodes().CreateVerificationCode(&verCode)
			if err != nil {
				c.JSON(http.StatusInternalServerError, loginResponse{
					StatusCode: http.StatusInternalServerError,
					Status:     "error",
					Message:    "An error occured",
				})
				return
			}

			// Make 6 digit 2FA code from first and last 3 digits of UUID Token
			verCodeStr := verCode.ID.String()
			code := verCodeStr[:3] + verCodeStr[len(verCodeStr)-3:]

			// Send as email
			verEmail := interfaces.EmailMsg{
				Subject: "Your FamTrust 2FA Code",
				From:    "FamTrust <biz@famtrust.biz>",
				To:      user.Email,
				BodyText: fmt.Sprintf("Hello there! \n"+
					"You've requested a 2FA code to login to your FamTrust account. \n"+
					"Use the code below to login. \n\n\n"+
					"\"%s\"", code),
			}

			if err = uh.mailer.SendMail(&verEmail); err != nil {
				c.JSON(http.StatusInternalServerError, loginResponse{
					StatusCode: http.StatusInternalServerError,
					Status:     "error",
					Message:    "Failed to user's 2FA Code, an error occured",
				})
				return

			} else {
				c.JSON(http.StatusOK, loginResponse{
					StatusCode: http.StatusOK,
					Status:     "success",
					Message:    "User has 2FA. Code has been sent to user's email",
				})
				return
			}

		} else {
			code, err := uh.models.VerCodes().Get2FACodeByUserID(user.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, loginResponse{
					StatusCode: http.StatusInternalServerError,
					Status:     "error",
					Message:    "An error occured",
				})
				return
			}

			codeStr := code.ID.String()
			if (codeStr[:3] + codeStr[len(codeStr)-3:]) != twoFACode {
				c.JSON(http.StatusUnauthorized, loginResponse{
					StatusCode: http.StatusUnauthorized,
					Status:     "error",
					Message:    "Invalid or Expired 2FA Code. Use the latest code",
				})
				return
			}

		}
	}

	token, err := jwtmod.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error has occured",
		})
		return
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
