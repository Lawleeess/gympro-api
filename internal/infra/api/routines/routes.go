package routines

import (
	"github.com/CValier/gympro-api/internal/infra/api/middlewares"
	"github.com/CValier/gympro-api/internal/infra/repositories/firebasedb"
	"github.com/CValier/gympro-api/internal/infra/repositories/firestoredb"
	"github.com/CValier/gympro-api/internal/pkg/service/auth"
	"github.com/CValier/gympro-api/internal/pkg/service/routine"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes adds user routes to the main engine router
func RegisterRoutes(e *gin.Engine) {
	routinesRoutes := e.Group("/api/v1/routines")

	repo := firestoredb.NewClient()
	authProvider := firebasedb.NewClient()
	authSvc := auth.NewAuthService(authProvider)
	routineService := routine.NewUserService(repo, authSvc)
	authHandler := newHandler(routineService)

	routinesRoutes.Use(middlewares.AuthenticateUser())

	routinesRoutes.POST("", authHandler.addRoutine)
	routinesRoutes.PUT("/image/:id", authHandler.updateRoutineUser)
	routinesRoutes.GET("", authHandler.getRoutines)

}
