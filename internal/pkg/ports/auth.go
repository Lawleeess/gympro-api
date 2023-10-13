package ports

import (
	"context"
	"mime/multipart"

	"firebase.google.com/go/auth"
	"github.com/CValier/gympro-api/internal/pkg/entity"
)

// AuthProvider is the signature to perform auth operations
type AuthProvider interface {
	SignInWithPass(ctx context.Context, credentials *entity.StandardLoginCredentials) (string, error)
	SignUpWithEmailAndPass(email, password string) (string, error) // Returns the id of the created user and error.
	SignInWithTokenClaims(ctx context.Context, token string) (*entity.SignWithCustomTokenResp, error)
	GenerateCustomToken(ctx context.Context, userID string, claims map[string]interface{}) (string, error)
	VerifyToken(token string) (*auth.Token, error)
	RevokeUserTokens(userID string) error
	RemoveUser(idToken string)
	UpdateUserImage(fileInput multipart.File, userID string) (string, error)
}
