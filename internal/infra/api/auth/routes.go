package auth

import (
	"github.com/CValier/gympro-api/internal/infra/repositories/firebasedb"
	"github.com/CValier/gympro-api/internal/infra/repositories/firestoredb"
	"github.com/CValier/gympro-api/internal/pkg/service/auth"
	"github.com/CValier/gympro-api/internal/pkg/service/user"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes adds auth routes to the main engine router
func RegisterRoutes(e *gin.Engine) {
	authRoutes := e.Group("/api/v1/auth")

	repo := firestoredb.NewClient()
	authProvider := firebasedb.NewClient()
	authSvc := auth.NewAuthService(authProvider)
	userService := user.NewUserService(repo, authSvc)
	authHandler := newHandler(userService)

	authRoutes.POST("/login", authHandler.signInWithPassword)
	authRoutes.POST("/signup", authHandler.signUp)
	authRoutes.POST("/verifyEmail", authHandler.VerifyOrRecoverEmail)
	authRoutes.POST("/sendOobCode", authHandler.VerifyOobCode)
	authRoutes.POST("/recover", authHandler.VerifyOrRecoverEmail)
}
