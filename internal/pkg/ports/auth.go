package ports

import (
	"context"
	"mime/multipart"

	"firebase.google.com/go/auth"
	"github.com/CValier/gympro-api/internal/pkg/entity"
)

// AuthProvider is the signature to perform auth operations
type AuthProvider interface {
	SignInWithPass(ctx context.Context, credentials *entity.StandardLoginCredentials) (*entity.SignWithCustomTokenResp, error)
	SignUpWithEmailAndPass(email, pass string) (*entity.SignWithCustomTokenResp, error)
	SignInWithTokenClaims(ctx context.Context, token string) (*entity.SignWithCustomTokenResp, error)
	GenerateCustomToken(ctx context.Context, userID string, claims map[string]interface{}) (string, error)
	VerifyToken(token string) (*auth.Token, error)
	RevokeUserTokens(userID string) error
	RemoveUser(idToken string) error
	UpdateUserImage(fileInput multipart.File, userID string) (string, error)
	UpdateRoutineImage(img multipart.File, id string) (string, error)

	VerifyOrRecoverEmail(ctx context.Context, creds *entity.UserRequestType) (string, error)
	VerifyOobCode(ctx context.Context, creds *entity.OobCode) (bool, error)

	ExchangeRefreshForIdToken(refreshToken string) *entity.SignWithCustomTokenResp
}
