package handlers

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/InternPulse/famtrust-backend-auth/internal/jwtmod"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandlers struct {
	models interfaces.Models
	mailer interfaces.Mailer
}

func (uh *UserHandlers) GetPermissions(permissions []interfaces.Permission) []string {
	var list []string
	for _, perm := range permissions {
		list = append(list, perm.ID)
	}
	return list
}

// @Summary		Login to FamTrust (Supports 2FA by Email)
// @Description	Login to FamTrust (Supports 2FA by Email)
// @Tags			User-Authentication
// @ID				login
// @Accept			json
// @Produce		json
// @Failure		401			{object}	loginSampleResponseError401
// @Failure		500			{object}	loginSampleResponseError500
// @Success		200			{object}	loginSampleResponse200
// @Param			Credentials	body		loginRequest	true	"User Credentials"
// @Param			2FACode		query		string			false	"User 2FA Code"
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

	var codeStr string

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

			codeStr = code.ID.String()
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

	if codeStr != "" {
		err := uh.models.VerCodes().Delete2FACodeByUserID(user.ID)
		if err != nil {
			log.Printf("Unable to delete 2FA Verifcation Code: %v", err)
		}
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

	permissions := uh.GetPermissions(user.Role.Permissions)

	// user payload
	role := role{
		Id:          user.Role.ID,
		Permissions: permissions,
	}
	userPayload := loginUserData{
		Id:           user.ID,
		Email:        user.Email,
		Has2FA:       user.Has2FA,
		DefaultGroup: user.DefaultGroup,
		IsVerified:   user.IsVerified,
		IsFreezed:    user.IsFreezed,
		LastLogin:    user.LastLogin,
		Role:         role,
	}

	payload := verifyResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "User session is valid",
		Token:      token.(string),
		User:       userPayload,
	}

	c.JSON(http.StatusOK, payload)
}

// @Summary		Create an Admin/Main User Account
// @Description	Create an Admin/Main User Account
// @Tags			User-Accounts
// @ID				signup
// @Accept			mpfd
// @Produce		json
// @Failure		400
// @Failure		500	{object}	loginSampleResponseError500
// @Success		201
// @Param			email		formData	string	true	"Email of the new user"
// @Param			password	formData	string	true	"Password of the new user"
// @Param			has2FA		formData	string	false	"Optional true or false value to set new user 2FA preference"
// @Router			/signup [post]
func (uh *UserHandlers) Signup(c *gin.Context) {
	var user interfaces.User

	switch {

	case strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded"):
	case strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data"):

		email := c.PostForm("email")
		password := c.PostForm("password")
		has2FAStr := c.PostForm("has2FA")

		if email == "" || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Incomplete signup Info",
			})
			return
		}

		if has2FAStr != "" {
			has2FA, err := strconv.ParseBool(has2FAStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "2FA value must be either true or false",
				})
				return
			}

			user.Has2FA = has2FA
		}

		// Generate user password
		bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"status":     "error",
				"message":    "Error parsing user password",
			})
			return
		}
		passwordHash := string(bytes)

		user.Email = email
		user.PasswordHash = passwordHash
		// Set admin as default role ID for user created via /signup
		user.RoleID = "admin"
		user.LastLogin = time.Now()

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"status":     "error",
			"message":    "Invalid signup data. Submit valid form-data",
			"error": gin.H{
				"error": fmt.Sprintf("You made use of a :%s: header", c.GetHeader("Content-Type")),
			},
		})
		return
	}

	err := uh.models.Users().CreateUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "A user with that email already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "error",
			"message":    "An error occured, failed to create User",
		})
		return
	}

	token, err := jwtmod.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "Error processing user sign in",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"status":     "success",
		"message":    "User created successfully. Proceed to verify email",
		"token":      token,
	})
}

// @Summary		Create a Sub-User/Member User Account
// @Description	Create a Sub-User/Member User Account - Requires the canCreateUsers permission
// @Tags			User-Accounts
// @ID				create-user
// @Security		BearerAuth
// @Accept			mpfd
// @Produce		json
// @Failure		400
// @Failure		500	{object}	loginSampleResponseError500
// @Success		201
// @Param			email		formData	string	true	"Email of the new user"
// @Param			password	formData	string	true	"Password of the new user"
// @Param			roleID		formData	string	false	"Optional Role ID string for new user. Defaults to 'member' if not specified"
// @Param			has2FA		formData	string	false	"Optional true or false value to set new user 2FA preference"
// @Router			/users [post]
func (uh *UserHandlers) CreateUser(c *gin.Context) {
	var user interfaces.User

	userWhoCreatesID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured",
		})
		return
	}

	// validate the user against the database
	userWhoCreates, err := uh.models.Users().GetUserByID(userWhoCreatesID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured, couldn't verify user",
		})
		return
	}

	// Confirm User 'canListUsers'
	permissions := uh.GetPermissions(userWhoCreates.Role.Permissions)
	if !slices.Contains(permissions, "canCreateUsers") {
		c.JSON(http.StatusUnauthorized, loginResponse{
			StatusCode: http.StatusUnauthorized,
			Status:     "error",
			Message:    "User does not have the necessary permission to perfom action",
		})
		return
	}

	switch {

	case strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded"):
	case strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data"):

		email := c.PostForm("email")
		password := c.PostForm("password")
		has2FAStr := c.PostForm("has2FA")
		roleID := c.PostForm("roleID")

		if email == "" || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Incomplete signup Info",
			})
			return
		}

		if has2FAStr != "" {
			has2FA, err := strconv.ParseBool(has2FAStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "2FA value must be either true or false",
				})
				return
			}

			user.Has2FA = has2FA
		}

		if roleID != "" {
			user.RoleID = roleID
		} else {
			user.RoleID = "member"
		}

		// Generate user password
		bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"status":     "error",
				"message":    "Error parsing user password",
			})
			return
		}
		passwordHash := string(bytes)

		user.DefaultGroup = userWhoCreates.DefaultGroup
		user.Email = email
		user.PasswordHash = passwordHash
		user.LastLogin = time.Now()

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"status":     "error",
			"message":    "Invalid signup data. Submit valid form-data",
			"error": gin.H{
				"error": fmt.Sprintf("You made use of a :%s: header", c.GetHeader("Content-Type")),
			},
		})
		return
	}

	err = uh.models.Users().CreateUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "A user with that email already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "error",
			"message":    "An error occured, failed to create User",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"status":     "success",
		"message":    "User created successfully. User can sign in now",
	})
}

// @Summary		Get All Users in Group
// @Description	Get All Users in Group - Requires the canListUsers permission
// @Tags			User-Accounts
// @ID				all-users-in-group
// @Security 		BearerAuth
// @Accept			json
// @Produce		json
// @Failure		400
// @Failure		500	{object}	loginSampleResponseError500
// @Success		201
// @Router			/users [get]
func (uh *UserHandlers) GetUsersByDefaultGroup(c *gin.Context) {

	UserID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured",
		})
		return
	}

	// validate the user against the database
	user, err := uh.models.Users().GetUserByID(UserID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured, couldn't verify user",
		})
		return
	}

	// Confirm User 'canListUsers'
	permissions := uh.GetPermissions(user.Role.Permissions)
	if !slices.Contains(permissions, "canListUsers") {
		c.JSON(http.StatusUnauthorized, loginResponse{
			StatusCode: http.StatusUnauthorized,
			Status:     "error",
			Message:    "User does not have the necessary permissions to perfom action",
		})
		return
	}

	// get the users from the database
	users, err := uh.models.Users().GetUsersByDefaultGroup(user.DefaultGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured while retrieving users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"status":     "success",
		"message":    "Users retrieved successfully",
		"users":      users,
	})
}

// @Summary		Get One User
// @Description	Get One User in User's Group - Requires the canListUsers permission
// @Tags			User-Accounts
// @ID				one-user-group
// @Security 		BearerAuth
// @Accept			json
// @Produce		json
// @Failure		400
// @Failure		401
// @Failure		500	{object}	loginSampleResponseError500
// @Success		201
// @Param			email		query	string	false	"User Email"
// @Param			code		query	string	false	"Password reset code"
// @Param			newPassword		formData	string	false	"New user password"
// @Router			/reset-password [get]
func (uh *UserHandlers) GetUserByDefaultGroup(c *gin.Context) {

	UserID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured",
		})
		return
	}

	userToGetStr := c.Param("userID")
	userToGetID, err := uuid.Parse(userToGetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, loginResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Invalid user ID",
		})
		return
	}

	// validate the user against the database
	user, err := uh.models.Users().GetUserByID(UserID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured, couldn't verify user",
		})
		return
	}

	// Confirm User 'canListUsers'
	permissions := uh.GetPermissions(user.Role.Permissions)
	if !slices.Contains(permissions, "canListUsers") {
		c.JSON(http.StatusUnauthorized, loginResponse{
			StatusCode: http.StatusUnauthorized,
			Status:     "error",
			Message:    "User does not have the necessary permissions to perfom action",
		})
		return
	}

	// get the users from the database
	userToGet, err := uh.models.Users().GetUserByDefaultGroup(userToGetID, user.DefaultGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured while retrieving user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"status":     "success",
		"message":    "User retrieved successfully",
		"user":       userToGet,
	})
}

func (uh *UserHandlers) ResetPassword(c *gin.Context) {
	resetCodeStr := c.Query("code")
	email := c.Query("email")
	newPass := c.PostForm("newPassword")

	if resetCodeStr == "" {
		if email == "" {
			c.JSON(http.StatusBadRequest, loginResponse{
				StatusCode: http.StatusBadRequest,
				Status:     "error",
				Message:    "User email or password reset code not provided",
			})
			return
		} else {
			user, err := uh.models.Users().GetUserByEmail("")
			if err != nil {
				c.JSON(http.StatusOK, loginResponse{
					StatusCode: http.StatusOK,
					Status:     "success",
					Message:    "Passord reset link sent if user exists",
				})
				return
			}

			resetCode := interfaces.VerCode{
				UserID: user.ID,
				Type:   "password",
			}

			err = uh.models.VerCodes().CreateVerificationCode(&resetCode)
			if err != nil {
				c.JSON(http.StatusInternalServerError, loginResponse{
					StatusCode: http.StatusInternalServerError,
					Status:     "error",
					Message:    "Failed to create user password reset link",
				})
				return
			}

			resetLink := "https://" + c.Request.Host + c.Request.URL.Path + "/reset-password/reset" + "?code=" + resetCode.ID.String()

			// Send as email
			verEmail := interfaces.EmailMsg{
				Subject: "Reset your password",
				From:    "FamTrust <biz@famtrust.biz>",
				To:      user.Email,
				BodyText: fmt.Sprintf("Hello there! \n"+
					"You've requested a password reset link for your FamTrust account. \n"+
					"Click on the link below to reset your password. \n\n\n"+
					"\"%s\"", resetLink),
			}

			if err = uh.mailer.SendMail(&verEmail); err != nil {
				c.JSON(http.StatusInternalServerError, loginResponse{
					StatusCode: http.StatusInternalServerError,
					Status:     "error",
					Message:    "Failed to send user password reset link",
				})
				return

			} else {
				c.JSON(http.StatusOK, loginResponse{
					StatusCode: http.StatusOK,
					Status:     "success",
					Message:    "Passord reset link sent if user exists",
				})
				return
			}
		}
	} else {
		if newPass == "" {
			c.JSON(http.StatusBadRequest, loginResponse{
				StatusCode: http.StatusBadRequest,
				Status:     "error",
				Message:    "No new password specified",
			})
			return
		}

		switch {

		case strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded"):
		case strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data"):
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Invalid data. Submit valid form-data",
				"error": gin.H{
					"error": fmt.Sprintf("You made use of a :%s: header", c.GetHeader("Content-Type")),
				},
			})
			return
		}

		resetCode, err := uuid.Parse(resetCodeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, loginResponse{
				StatusCode: http.StatusBadRequest,
				Status:     "error",
				Message:    "Invalid reset token",
			})
			return
		}

		code, err := uh.models.VerCodes().GetResetCodeByID(resetCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, loginResponse{
				StatusCode: http.StatusBadRequest,
				Status:     "error",
				Message:    "Invalid reset token",
			})
			return
		}

		// Generate new user password
		bytes, err := bcrypt.GenerateFromPassword([]byte(newPass), 14)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"status":     "error",
				"message":    "Error parsing new user password",
			})
			return
		}
		passwordHash := string(bytes)

		user, err := uh.models.Users().GetUserByID(code.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"status":     "error",
				"message":    "Error retrieving user",
			})
			return
		}

		user.PasswordHash = passwordHash

		err = uh.models.Users().UpdateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"status":     "error",
				"message":    "Error updating user with new password",
			})
			return
		}

		if err := uh.models.VerCodes().DeleteResetCodeByUserID(code.UserID); err != nil {
			log.Printf("Failed to delete Email verification code: %v", err)
		}
	}

	c.JSON(http.StatusOK, loginResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Password successfully updated",
	})
}
