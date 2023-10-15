package user

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
	authRoutes := e.Group("/api/v1/users")

	repo := firestoredb.NewClient()
	authProvider := firebasedb.NewClient()
	authSvc := auth.NewAuthService(authProvider)
	userService := user.NewUserService(repo, authSvc)
	authHandler := newHandler(userService)

	authRoutes.Use(middlewares.AuthenticateUser())

	authRoutes.PUT("/image/:user_id", authHandler.updateImageUser)
	authRoutes.PUT("/:user_id", authHandler.updateUser)
}
