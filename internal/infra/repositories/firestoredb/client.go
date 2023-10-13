package firestoredb

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/CValier/gympro-api/internal/pkg/config"
	"github.com/CValier/gympro-api/internal/pkg/entity"
	"github.com/epa-datos/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type firestoreClient struct {
	client *firestore.Client
}

var (
	onceFirestoreClient sync.Once
	firestoreClie       *firestoreClient
)

// NewClient returns a new instance for firestore client
func NewClient() *firestoreClient {
	onceFirestoreClient.Do(func() {
		ctx := context.Background()
		opt := option.WithCredentialsJSON([]byte(config.CfgIn.ServiceCredentialJSON))
		client, err := firestore.NewClient(ctx, config.CfgIn.GoogleProjectID, opt)
		if err != nil {
			logrus.Error("firestoredb.NewClient failed to create a new client: " + err.Error())
		}
		firestoreClie = &firestoreClient{
			client: client,
		}

	})
	return firestoreClie
}

func (f *firestoreClient) AddUser(user *entity.User) error {
	if user.ID == "" {
		return errors.Build(
			errors.Operation("firestoredb.AddUser"),
			errors.InternalError,
			errors.Message("User id can't be empty"),
		)
	}
	_, err := f.client.Collection("users").Doc(user.ID).Set(context.Background(), user)
	return err
}

// GetUserByEmail gets a user with the given email
func (f *firestoreClient) GetUserByEmail(email string) (*entity.User, error) {
	user := &entity.User{}

	usersCollection := f.client.Collection("users")

	// Creating the query
	query := usersCollection.Where("email", "==", email)

	iter := query.Documents(context.Background())
	defer iter.Stop()

	// Reading the response query
	doc, err := iter.Next()
	if err != nil {
		return nil, errors.Build(
			errors.NotFound,
			errors.Message(err.Error()),
		)
	}

	if err := doc.DataTo(user); err != nil {
		logrus.Error("firestoredb.GetUserByEmail, Failed to decode user: " + err.Error())
		return nil, errors.Build(
			errors.InternalError,
		)
	}

	user.ID = doc.Ref.ID

	return user, nil
}

// GetUserByID returns a user from firestore, according to given user id.
func (f *firestoreClient) GetUserByID(userID string) (*entity.User, error) {
	scope := errors.Operation("firestoredb.GetUserByID")

	doc, err := f.client.Collection("users").Doc(userID).Get(context.Background())
	// If is not reference to the document. Then the user with given ID is not found.
	if err != nil {
		return nil, errors.Build(
			scope,
			errors.NotFound,
			errors.Message("User not found: "+userID),
		)
	}

	// Deserealize document into proper entity.
	user := new(entity.User)
	if err := doc.DataTo(user); err != nil {
		return nil, errors.Build(
			scope,
			errors.InternalError,
			errors.Message("Failed to deserealize doc: "+err.Error()),
		)
	}
	user.ID = doc.Ref.ID

	return user, nil
}

// GetAllUsersCount returns an int that means the total number of documents on the firestore collection
func (f *firestoreClient) GetAllUsersCount() (int, error) {
	query := f.client.Collection("users")
	docs, err := query.Documents(context.Background()).GetAll()
	return len(docs), err
}

// GetAllUsers returns documents from the users collection between the offset and limit params.
func (f *firestoreClient) GetUsersByPage(offset, limit int64, department, filter string) ([]*entity.User, error) {
	var users []*entity.User
	query := f.client.Collection("users").OrderBy("name", firestore.Asc).Offset(int(offset)).Limit(int(limit))
	iter := query.Documents(context.Background())
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Build(
				errors.Operation("firestoredb.GetAllUsers"),
				errors.InternalError,
				errors.Message("Failed to get documents: "+err.Error()),
			)
		}
		user := new(entity.User)
		if err := doc.DataTo(user); err != nil {
			return nil, errors.Build(
				errors.Operation("firestoredb.GetAllUsers"),
				errors.InternalError,
				errors.Message("Failed to desirealize obj: "+err.Error()),
			)
		}
		user.ID = doc.Ref.ID
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser update a user from firestore, according to given user id.
func (f *firestoreClient) UpdateUser(userID string, user *entity.User) error {
	_, err := f.client.Collection("users").
		Doc(userID).
		Set(context.Background(), user)
	return err
}

// UpdateUser update a user from firestore, according to given user id.
func (f *firestoreClient) UpdateImageUser(userID string, url string) error {
	userDoc := f.client.Collection("users").Doc(userID)

	// Making the update over firestore collection.
	_, err := userDoc.Update(
		context.Background(),
		[]firestore.Update{
			{Path: "url_image", Value: url},
		},
	)
	return err
}
