package ports

import (
	"mime/multipart"

	"firebase.google.com/go/auth"
)

// FirebaseCli is signature to perform operations over firebase
type FirebaseCli interface {
	GenerateCustomToken(userID string, claims map[string]interface{}) (string, error)
	VerifyToken(token string) (*auth.Token, error)
	RevokeUserTokens(userID string) error
	UpdateUserImage(fileInput multipart.File, userID string) (string, error)
	UpdateRoutineImage(fileInput multipart.File, id string) (string, error)
}
