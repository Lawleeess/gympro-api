package api

import (
	"github.com/CValier/gympro-api/internal/infra/api/auth"
	"github.com/CValier/gympro-api/internal/infra/api/routines"
	"github.com/CValier/gympro-api/internal/infra/api/user"
	usrManagement "github.com/CValier/gympro-api/internal/infra/api/user-managment"

	"github.com/gin-gonic/gin"
)

func registerRoutes(e *gin.Engine) {

	auth.RegisterRoutes(e)

	user.RegisterRoutes(e)

	usrManagement.RegisterRoutes(e)

	routines.RegisterRoutes(e)
}
