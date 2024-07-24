package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary		Retrieve User Profile Details
// @Description	Retrieve User Profile Details
// @Tags			User-Profiles
// @ID				profile
// @Accept			json
// @Produce		json
// @Failure		401	{object}	loginSampleResponseError401
// @Failure		500	{object}	loginSampleResponseError500
// @Success		200	{object}	profileSampleResponse200
// @Security		BearerAuth
// @Router			/profile [get]
func (uh *UserHandlers) GetUserProfileByID(c *gin.Context) {
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
	profile, err := uh.models.Users().GetUserProfileByID(UserID.(uuid.UUID))
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			c.JSON(http.StatusNotFound, loginResponse{
				StatusCode: http.StatusNotFound,
				Status:     "error",
				Message:    "User doesn't have a profile",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, loginResponse{
				StatusCode: http.StatusInternalServerError,
				Status:     "error",
				Message:    "An error occured",
			})
			return
		}

	}

	payload := gin.H{
		"statusCode": http.StatusOK,
		"status":     "success",
		"message":    "User profile retrieved successfully",
		"profile":    profile,
	}

	c.JSON(http.StatusOK, payload)
}
