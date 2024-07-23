package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VerificationHandlers struct {
	models interfaces.Models
	mailer interfaces.Mailer
}

// @Summary		Send User-Email Verification Token
// @Description	Send User-Email Verification Token
// @Tags			Verifications
// @ID				send-verify-token
// @Security BearerAuth
// @Produce		json
// @Failure		400
// @Failure		500
// @Success		200
// @Router			/verify-email [get]
func (v *VerificationHandlers) VerifyEmail(c *gin.Context) {

	UserID, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured",
		})
		return
	}

	// retrieve the user from the database
	user, err := v.models.Users().GetUserByID(UserID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "An error occured",
		})
		return
	}

	if user.IsVerified {
		c.JSON(http.StatusBadRequest, loginResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "User is already verified",
		})
		return

	} else {

		verCode := interfaces.VerCode{
			UserID: user.ID,
		}

		err := v.models.VerCodes().CreateVerificationCode(&verCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, loginResponse{
				StatusCode: http.StatusInternalServerError,
				Status:     "error",
				Message:    "An error occured",
			})
			return
		}

		// create verification link
		verLink := c.Request.Host + c.Request.URL.Path + "/verify" + "?code=" + verCode.ID.String()

		// send as email
		verEmail := interfaces.EmailMsg{
			Subject:  "Verify your FamTrust Email",
			From:     "biz@famtrust.biz",
			To:       user.Email,
			BodyText: fmt.Sprintf("Hello there! \nWelcome to FamTrust. \nClick the link below to verify your Email Address. \n\n\n %s", verLink),
		}

		if err = v.mailer.SendMail(&verEmail); err != nil {
			c.JSON(http.StatusInternalServerError, loginResponse{
				StatusCode: http.StatusInternalServerError,
				Status:     "error",
				Message:    "Failed to send verification email, an error occured",
			})
			return
		} else {
			c.JSON(http.StatusOK, loginResponse{
				StatusCode: http.StatusOK,
				Status:     "success",
				Message:    "Verification email sent successfully",
			})
			return
		}

	}
}

// @Summary		Verify User Email Address via Token
// @Description	Verify User Email Address via Token
// @Tags			Verifications
// @ID				verify-email-token
// @Produce		json
// @Failure		400
// @Failure		500
// @Success		200
// @Param code query string true "Email verification Token"
// @Router			/verify-email/verify [get]
func (v *VerificationHandlers) VerifyEmailToken(c *gin.Context) {
	verCode := c.Query("code")

	if verCode == "" {
		c.JSON(http.StatusBadRequest, loginResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Invalid verification code",
		})
		return
	}

	codeID, err := uuid.Parse(verCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, loginResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Invalid verification code",
		})
		return
	}

	code, err := v.models.VerCodes().GetCodeByID(codeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, loginResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Invalid verification code",
		})
		return
	}

	if err = v.models.Users().SetIsVerified(code.UserID, true); err != nil {
		c.JSON(http.StatusInternalServerError, loginResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "Failed to update user status to verified",
		})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "Email successfully verified",
	})
}

// @Summary		Verify User Signup NIN
// @Description	Verify User Signup NIN
// @Tags			Verifications
// @ID				verify-nin
// @Produce		json
// @Failure		400
// @Failure		500
// @Success		200
// @Param nin query int true "NIN"
// @Router			/verify-nin [get]
func (v *VerificationHandlers) VerifyNIN(c *gin.Context) {
	ninStr := c.Query("nin")
	if ninStr != "" {
		nin, err := strconv.Atoi(ninStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Invalid NIN",
			})
			return
		}
		if len(ninStr) != 10 {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Invalid NIN",
			})
			return
		}
		_, err = v.models.Users().GetUserByNIN(nin)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": http.StatusOK,
				"status":     "success",
				"message":    "NIN is valid and un-used",
			})
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"status":     "error",
				"message":    "A user with this NIN already exists",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"status":     "error",
			"message":    "You must specify an NIN",
		})
		return
	}
}

// @Summary		Verify User Signup BVN
// @Description	Verify User Signup BVN
// @Tags			Verifications
// @ID				verify-bvn
// @Produce		json
// @Failure		400
// @Failure		500
// @Success		200
// @Param bvn query int true "BVN"
// @Router			/verify-bvn [get]
func (v *VerificationHandlers) VerifyBVN(c *gin.Context) {
	bvnStr := c.Query("bvn")
	if bvnStr != "" {
		bvn, err := strconv.Atoi(bvnStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Invalid BVN",
			})
			return
		}
		if len(bvnStr) != 10 {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Invalid BVN",
			})
			return
		}
		_, err = v.models.Users().GetUserByBVN(bvn)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"statusCode": http.StatusOK,
				"status":     "success",
				"message":    "BVN is valid and un-used",
			})
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"status":     "error",
				"message":    "A user with this BVN already exists",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"status":     "error",
			"message":    "You must specify an BVN",
		})
		return
	}
}
