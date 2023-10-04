package auth

import (
	"net/http"

	"github.com/CValier/gympro-api/internal/pkg/entity"
	"github.com/CValier/gympro-api/internal/pkg/ports"
	"github.com/epa-datos/errors"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	userService ports.UserService
}

func newHandler(service ports.UserService) *authHandler {
	return &authHandler{
		userService: service,
	}
}

func (u *authHandler) signUp(c *gin.Context) {
	user := &entity.User{}

	if err := c.Bind(user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := u.userService.CreateUser(c, user); err != nil {
		errors.JSON(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
func (u *authHandler) signInWithPassword(c *gin.Context) {
	var credentials entity.StandardLoginCredentials

	if err := c.Bind(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, errors.Build(
			errors.Message("Failed to bind credentials: "+err.Error()),
		))
		return
	}

	authResponse, err := u.userService.SignInWithPass(c, &credentials)
	if err != nil {
		errors.JSON(c, err)
		return
	}

	c.JSON(http.StatusOK, authResponse)
}
