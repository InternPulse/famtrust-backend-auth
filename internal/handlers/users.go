package handlers

import (
	"fmt"
	"log"
	"net/http"
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
		User:       userPayload,
	}

	c.JSON(http.StatusOK, payload)
}

// @Summary		Create an Admin/Main User Account
// @Description	Create an Admin/Main User Account
// @Tags			User-Accounts
// @ID				signup
// @Accept			json
// @Produce		json
// @Failure		400
// @Failure		500	{object}	loginSampleResponseError500
// @Success		201
// @Param email formData string true "Email of the new user"
// @Param password formData string true "Password of the new user"
// @Param has2FA formData string false "Optional true or false value to set new user 2FA preference"
// @Router			/signup [post]
func (uh *UserHandlers) Signup(c *gin.Context) {
	var user interfaces.User

	switch {

	case strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded"):
	case strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data"):
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Failed to parse form",
			})
			return
		}

		email := form.Value["email"][0]
		password := form.Value["password"][0]
		has2FAStr := form.Value["has2FA"][0]

		// image := form.File["image"][0]
		// if image != nil {
		// 	c.SaveUploadedFile(image, fmt.Sprintf("%s/%s", uploadDir, image.Filename))
		// }

		// price, err := strconv.ParseFloat(priceStr[0], 64)
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"status":  "error",
		// 		"message": "Invalid price, price must be a number",
		// 	})
		// 	return
		// }

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

		user = interfaces.User{
			Email:        email,
			PasswordHash: passwordHash,
			RoleID:       "admin",
			LastLogin:    time.Now(),
		}

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
// @Description	Create a Sub-User/Member User Account
// @Tags			User-Accounts
// @ID				create-user
// @Accept			json
// @Produce		json
// @Failure		400
// @Failure		500	{object}	loginSampleResponseError500
// @Success		201
// @Param email formData string true "Email of the new user"
// @Param password formData string true "Password of the new user"
// @Param roleID formData string false "Optional Role ID string for new user. Defaults to 'member' if not specified"
// @Param has2FA formData string false "Optional true or false value to set new user 2FA preference"
// @Router			/create-user [post]
func (uh *UserHandlers) CreateUser(c *gin.Context) {
	var user interfaces.User

	switch {

	case strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded"):
	case strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data"):
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Failed to parse form",
			})
			return
		}

		email := form.Value["email"][0]
		password := form.Value["password"][0]
		has2FAStr := form.Value["has2FA"][0]
		roleID := form.Value["roleID"][0]

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

		user = interfaces.User{
			Email:        email,
			PasswordHash: passwordHash,
			LastLogin:    time.Now(),
		}

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
