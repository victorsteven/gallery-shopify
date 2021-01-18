package handler

import (
	"fmt"
	"gallery/server/domain"
	"gallery/server/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type service struct {
	imgDomain   domain.ImageService
	userService domain.UserService
}

func NewHandlerService(imgDomain domain.ImageService, userService domain.UserService) *service {
	return &service{imgDomain, userService}
}

func (s *service) DeleteImage(c *gin.Context) {

	uid, err := utils.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	requestId := c.Param("imageId")

	imageId, err := strconv.Atoi(requestId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message":   "Invalid id provided",
		})
		return
	}

	// get the user by id:
	user, err := s.userService.GetUserByID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"message":   err.Error(),
		})
		return
	}

	img, err := s.imgDomain.GetImageByID(imageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"body":   err.Error(),
		})
		return
	}

	if user.ID != img.UserID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"message": fmt.Sprintf("Ensure that Image with the id %s exists and you have permission to delete it", imageId),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message":   "Success",
	})
}

func (s *service) BulkDeleteImage(c *gin.Context) {

	uid, err := utils.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	var request *domain.BulkImage

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "ensure that the image ids are valid",
		})
		return
	}

	// get the user by id:
	_, err = s.userService.GetUserByID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"message":   err.Error(),
		})
		return
	}

	err = s.imgDomain.BulkDeleteImage(request.ImageIds, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"message":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message":   "Successfully deleted images",
	})
}

func (s *service) Login(c *gin.Context) {

	var request *domain.User

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message":   "please provide valid inputs",
		})
		return
	}

	if err := utils.ValidateInputs(*request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"message":   err.Error(),
		})
		return
	}

	user, err := s.userService.GetUserByEmail(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"message":   err.Error(),
		})
		return
	}

	// check the password:
	err = utils.VerifyPassword(user.Password, request.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"message":   err.Error(),
		})
		return
	}

	token, err := utils.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"message":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"token":  token,
	})
}
