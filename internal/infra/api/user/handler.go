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

func (u *userHandler) saveProgressGoals(c *gin.Context) {
	userProgress := &entity.UserProgress{}

	if err := c.Bind(&userProgress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid format: " + err.Error(),
		})
		return
	}
	goals, errUpdate := u.userService.SaveUserProgress(c.Param("user_id"), userProgress)
	if errUpdate != nil {
		errors.JSON(c, errUpdate)
		return
	}

	c.JSON(http.StatusOK, goals)
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

func (u *userHandler) saveRoutines(c *gin.Context) {
	userRoutine := &entity.UserRoutine{}

	if err := c.Bind(&userRoutine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid format: " + err.Error(),
		})
		return
	}

	err := u.userService.AddRoutineToUser(c.Param("user_id"), userRoutine)
	if err != nil {
		errors.JSON(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
