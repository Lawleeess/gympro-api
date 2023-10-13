package user

import (
	"net/http"

	"github.com/CValier/gympro-api/internal/pkg/entity"
	"github.com/CValier/gympro-api/internal/pkg/ports"
	"github.com/epa-datos/errors"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService ports.UserService
}

func newHandler(service ports.UserService) *userHandler {
	return &userHandler{
		userService: service,
	}
}

func (u *userHandler) updateUser(c *gin.Context) {
	user := &entity.User{}

	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid format: " + err.Error(),
		})
		return
	}
	errUpdate := u.userService.UpdateUser(c.Param("user_id"), user)
	if errUpdate != nil {
		errors.JSON(c, errUpdate)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (u *userHandler) updateImageUser(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.Build(
			errors.Message("Failed to get image: "+err.Error()),
		))
		return
	}

	errUpdate := u.userService.UpdateImageUser(file, c.Param("user_id"))
	if errUpdate != nil {
		errors.JSON(c, errUpdate)
		return
	}

	c.JSON(http.StatusOK, nil)
}
