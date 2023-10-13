package firebasedb

import (
	"context"
	"io"
	"mime/multipart"
	"sync"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/CValier/gympro-api/internal/pkg/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type firebaseClient struct {
	client *auth.Client
}

var (
	onceFirebaseClient sync.Once
	firebaseClie       *firebaseClient
)

// NewClient returns a new instance for firebase client
func NewClient() *firebaseClient {
	onceFirebaseClient.Do(func() {
		ctx := context.Background()
		opt := option.WithCredentialsJSON([]byte(config.CfgIn.ServiceCredentialJSON))
		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			logrus.Error("firebasedb.NewClient failed to create a app: " + err.Error())
		}
		client, err := app.Auth(ctx)
		if err != nil {
			logrus.Error("firebasedb.NewClient failed to create a new client: " + err.Error())
		}
		firebaseClie = &firebaseClient{
			client: client,
		}
	})
	return firebaseClie
}

// GenerateCustomToken uses the SDK provided to create a token with claims encrypted
func (f *firebaseClient) GenerateCustomToken(userID string, claims map[string]interface{}) (string, error) {
	return f.client.CustomTokenWithClaims(context.Background(), userID, claims)
}

// VerifyToken uses the Firebase SDK to verify if a JWT is still valid.
func (f *firebaseClient) VerifyToken(token string) (*auth.Token, error) {
	return f.client.VerifyIDToken(context.Background(), token)
}

// RevokeUserTokens removes all active refresh tokens for a particular user.
func (f *firebaseClient) RevokeUserTokens(userID string) error {
	return f.client.RevokeRefreshTokens(context.Background(), userID)
}

func (f *firebaseClient) UpdateUserImage(fileInput multipart.File, userID string) (string, error) {

	path := "https://firebasestorage.googleapis.com/v0/b/gympro-400622.appspot.com/o/users%2F"

	conf := &firebase.Config{
		ProjectID:     config.CfgIn.GoogleProjectID,
		StorageBucket: config.CfgIn.GoogleProjectID + ".appspot.com",
	}
	opt := option.WithCredentialsJSON([]byte(config.CfgIn.ServiceCredentialJSON))

	app, err := firebase.NewApp(context.Background(), conf,
		opt)

	if err != nil {
		return "", nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	client, err := app.Storage(ctx)
	if err != nil {
		return "", nil
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return "", nil
	}
	idImg := uuid.New()
	nameImage := "profile_" + userID + ".png"

	object := bucket.Object("users/" + nameImage)
	writer := object.NewWriter(ctx)
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": idImg.String()}
	writer.ObjectAttrs.ContentType = "image/png"

	defer writer.Close()

	if _, err := io.Copy(writer, fileInput); err != nil {
		return "", nil
	}

	path += nameImage + "?alt=media&token=" + idImg.String()

	return path, nil
}
