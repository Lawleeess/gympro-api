package adminmanagment

import (
	"github.com/CValier/gympro-api/internal/infra/api/middlewares"
	"github.com/CValier/gympro-api/internal/infra/repositories/firebasedb"
	"github.com/CValier/gympro-api/internal/infra/repositories/firestoredb"
	"github.com/CValier/gympro-api/internal/pkg/service/auth"
	"github.com/CValier/gympro-api/internal/pkg/service/user"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes adds user routes to the main engine router
func RegisterRoutes(e *gin.Engine) {
	adminMngRoutes := e.Group("/api/v1/user-management")

	repo := firestoredb.NewClient()
	authProvider := firebasedb.NewClient()
	authSvc := auth.NewAuthService(authProvider)
	userService := user.NewUserService(repo, authSvc)
	authHandler := newHandler(userService)

	adminMngRoutes.Use(middlewares.AuthenticateUser())

	adminMngRoutes.GET("/users", authHandler.getAllUsers)
}
