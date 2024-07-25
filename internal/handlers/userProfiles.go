package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/InternPulse/famtrust-backend-auth/internal/interfaces"
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

// @Summary		Create New User Profile
// @Description	Create New User Profile
// @Tags			User-Profiles
// @ID				create-profile
// @Security		BearerAuth
// @Accept			multipart/form-data
// @Produce		json
// @Param			firstName		formData	string	true	"User's first name"
// @Param			lastName		formData	string	true	"User's last name"
// @Param			bio				formData	string	true	"User's biography"
// @Param			nin				formData	int		false	"User's National Identification Number"
// @Param			bvn				formData	int		false	"User's Bank Verification Number"
// @Param			profilePicture	formData	file	true	"User's profile picture"
// @Router			/profile/create [post]
func (uh *UserHandlers) CreateUserProfile(c *gin.Context) {

	var profile interfaces.UserProfile

	switch {

	case strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded"):
	case strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data"):

		UserID, exists := c.Get("UserID")
		if !exists {
			c.JSON(http.StatusInternalServerError, loginResponse{
				StatusCode: http.StatusInternalServerError,
				Status:     "error",
				Message:    "An error occured",
			})
			return
		}

		firstName := c.PostForm("firstName")
		lastName := c.PostForm("lastName")
		bio := c.PostForm("bio")
		ninStr := c.PostForm("nin")
		bvnStr := c.PostForm("bvn")

		if firstName == "" || lastName == "" || bio == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Incomplete required signup Info",
			})
			return
		}

		// Parse NIN
		if ninStr != "" {
			nin, err := strconv.ParseUint(ninStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Invalid value for NIN",
				})
				return
			}
			profile.NIN = uint(nin)
		}

		// Parse BVN
		if bvnStr != "" {
			bvn, err := strconv.ParseUint(bvnStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Invalid value for BVN",
				})
				return
			}
			profile.BVN = uint(bvn)
		}

		// Get profile picture
		profilePicture, err := c.FormFile("profilePicture")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Error parsing profile picture. Check your upload",
			})
			return
		}

		const maxUploadSize = 16 << 20
		if profilePicture.Size > maxUploadSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Picture file is too large",
			})
			return
		}

		// Generate a unique filename
		filename := filepath.Base(profilePicture.Filename)
		dst := filepath.Join("images", "profilePics", filename)
		dstURL := filepath.Join("images", "profile-pic", filename)

		// Save the file
		if err := c.SaveUploadedFile(profilePicture, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// User values
		profile.UserID = UserID.(uuid.UUID)
		profile.FirstName = firstName
		profile.LastName = lastName
		profile.Bio = bio
		profile.ProfilePictureUrl = dstURL

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

	err := uh.models.Users().CreateUserProfile(&profile)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "User already has a profile",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "error",
			"message":    "An error occured, failed to create user profile",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"status":     "success",
		"message":    "User profile created successfully",
	})
}

// @Summary			Update User Profile
// @Description		Update User Profile
// @Tags			User-Profiles
// @ID				update-profile
// @Security		BearerAuth
// @Accept			multipart/form-data
// @Produce			json
// @Param			firstName		formData	string	false	"User's first name"
// @Param			lastName		formData	string	false	"User's last name"
// @Param			bio				formData	string	false	"User's biography"
// @Param			nin				formData	int		false	"User's National Identification Number"
// @Param			bvn				formData	int		false	"User's Bank Verification Number"
// @Param			profilePicture	formData	file	false	"User's profile picture"
// @Router			/profile/update [put]
func (uh *UserHandlers) UpdateUserProfile(c *gin.Context) {

	var profile interfaces.UserProfile

	switch {

	case strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded"):
	case strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data"):

		UserID, exists := c.Get("UserID")
		if !exists {
			c.JSON(http.StatusInternalServerError, loginResponse{
				StatusCode: http.StatusInternalServerError,
				Status:     "error",
				Message:    "An error occured",
			})
			return
		}

		firstName := c.PostForm("firstName")
		lastName := c.PostForm("lastName")
		bio := c.PostForm("bio")
		ninStr := c.PostForm("nin")
		bvnStr := c.PostForm("bvn")

		if firstName != "" {
			profile.FirstName = firstName
		}

		if lastName != "" {
			profile.LastName = lastName
		}

		if bio != "" {
			profile.Bio = bio
		}

		// Parse NIN
		if ninStr != "" {
			nin, err := strconv.ParseUint(ninStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Invalid value for NIN",
				})
				return
			}
			profile.NIN = uint(nin)
		}

		// Parse BVN
		if bvnStr != "" {
			bvn, err := strconv.ParseUint(bvnStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": "Invalid value for BVN",
				})
				return
			}
			profile.BVN = uint(bvn)
		}

		// Check if a profile picture was submitted
		if fileHeader, err := c.FormFile("profilePicture"); err == nil && fileHeader != nil {
			// A file was submitted, process it
			profilePicture, err := c.FormFile("profilePicture")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"statusCode": http.StatusBadRequest,
					"status":     "error",
					"message":    "Error parsing profile picture. Check your upload",
				})
				return
			}

			const maxUploadSize = 16 << 20
			if profilePicture.Size > maxUploadSize {
				c.JSON(http.StatusBadRequest, gin.H{
					"statusCode": http.StatusBadRequest,
					"status":     "error",
					"message":    "Picture file is too large",
				})
				return
			}

			// Generate a unique filename
			filename := filepath.Base(profilePicture.Filename)
			dst := filepath.Join("images", "profilePics", filename)
			dstURL := filepath.Join("images", "profile-pic", filename)

			// Save the file
			if err := c.SaveUploadedFile(profilePicture, dst); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			profile.ProfilePictureUrl = dstURL

		} else if err != http.ErrMissingFile {
			// An error occurred while checking for the file (other than the file being missing)
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"status":     "error",
				"message":    "Error parsing image. Check your upload",
			})
			return
		}

		profile.UserID = UserID.(uuid.UUID)

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

	err := uh.models.Users().UpdateUserProfile(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"status":     "error",
			"message":    "An error occured, failed to update user profile",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"status":     "success",
		"message":    "User profile updated successfully",
	})
}

// @Summary		Get User Profile Picture
// @Description	Get User Profile Picture
// @Tags			User-Profiles
// @ID				get-profile-pic
// @Produce		json
// @Failure		404
// @Success		200
// @Param			imageName	path	string	true	"Picture Filename"
// @Router			/images/profile-pic/{imageName} [get]
func (uh *UserHandlers) GetProfilePicture(c *gin.Context) {
	imageName := c.Param("imageName")

	imagePath := filepath.Join("images", "profilePics", imageName)

	// Check if the file exists and is not a directory
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"status":     "error",
			"message":    "File does not exist",
		})
		return
	}

	c.File(imagePath)
}
