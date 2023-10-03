package api

import (
	"time"

	"github.com/CValier/gympro-api/internal/pkg/config"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

// RunServer initialize api server
func RunServer() {
	server := gin.Default()

	server.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET,POST,DELETE,PUT",
		RequestHeaders:  "Origin, Authorization, Content-Type, Access-Control-Allow-Origin",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	registerRoutes(server)

	_ = server.Run(":" + config.CfgIn.ServerPort)
}
