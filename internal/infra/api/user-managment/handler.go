package adminmanagment

import (
	"net/http"

	"github.com/CValier/gympro-api/internal/pkg/ports"
	"github.com/epa-datos/errors"
	"github.com/gin-gonic/gin"
)

type adminManagmentHandler struct {
	userService ports.UserService
}

func newHandler(service ports.UserService) *adminManagmentHandler {
	return &adminManagmentHandler{
		userService: service,
	}
}

func (u *adminManagmentHandler) getAllUsers(c *gin.Context) {
	c.Set("offset", c.Query("offset"))
	c.Set("limit", c.Query("limit"))
	c.Set("user_role", c.Query("user_role"))
	c.Set("filter", c.Query("filter"))

	users, err := u.userService.GetUsers(c)
	if err != nil {
		errors.JSON(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}
