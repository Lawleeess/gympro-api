package ports

import "firebase.google.com/go/auth"

// FirebaseCli is signature to perform operations over firebase
type FirebaseCli interface {
	GenerateCustomToken(userID string, claims map[string]interface{}) (string, error)
	VerifyToken(token string) (*auth.Token, error)
	RevokeUserTokens(userID string) error
}
