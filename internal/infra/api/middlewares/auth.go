package middlewares

import (
	"net/http"
	"strings"

	"github.com/CValier/gympro-api/internal/infra/repositories/firebasedb"
	"github.com/CValier/gympro-api/internal/infra/repositories/firestoredb"
	"github.com/CValier/gympro-api/internal/pkg/service/auth"
	"github.com/CValier/gympro-api/internal/pkg/service/user"
	"github.com/gin-gonic/gin"
)

// AuthenticateUser verifies if the user making the current request
// Is authenticated and his session is valid.
func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Init dependencies for user service
		userRepo := firestoredb.NewClient()
		authProvider := firebasedb.NewClient()
		authSvc := auth.NewAuthService(authProvider)

		userService := user.NewUserService(userRepo, authSvc)

		// Getting the JWT from the Authorization header.
		authorization := c.Request.Header.Get("Authorization")
		jwtToken := strings.TrimPrefix(authorization, "Bearer ")

		// Verifiying the token.
		userClaims, err := userService.VerifyToken(jwtToken)
		if err != nil {
			if err := c.AbortWithError(http.StatusUnauthorized, err); err != nil {
				c.JSON(http.StatusUnauthorized, nil)
			}
			return
		}

		// Setting user's info in the context to avoid
		// Making request to user repository every request
		c.Set("userID", userClaims["user_id"])
		c.Set("userEmail", userClaims["email"])
		c.Set("subscription", userClaims["subscription"])
		c.Set("modulesWithPermission", userClaims["modulesWithPermission"])
		c.Set("fullName", userClaims["fullName"])
		c.Set("birthday", userClaims["birthday"])
		c.Set("phone_number", userClaims["phone_number"])

		c.Next()
	}
}
