package user

import (
	"context"
	"fmt"
	"net/mail"
	"strconv"
	"time"

	"github.com/CValier/gympro-api/internal/pkg/entity"
	"github.com/CValier/gympro-api/internal/pkg/ports"
	"github.com/CValier/gympro-api/internal/pkg/utils"
	"github.com/epa-datos/errors"
	"github.com/sirupsen/logrus"
)

type userSvc struct {
	repo    ports.UsersRepository
	authSvc ports.AuthProvider
}

// NewUserService returns an instance for user service
func NewUserService(repo ports.UsersRepository, firebAuth ports.AuthProvider) *userSvc {
	return &userSvc{
		repo:    repo,
		authSvc: firebAuth,
	}
}

// CreateUser adds a new user to auth provider and users repository.
func (u *userSvc) CreateUser(ctx context.Context, user *entity.User) error {
	// 1. Check if the given email is a valid direction.
	if _, err := mail.ParseAddress(user.Email); err != nil {
		return err
	}

	// 2. Create user in our auth provider
	userID, err := u.authSvc.SignUpWithEmailAndPass(user.Email, user.Password)
	if err != nil {
		return err
	}

	user.ID = userID

	if user.Subscription == "" {
		timeNow := time.Now().Format("2006-01-02")
		// timeNow := strconv.Itoa(time.Now().Year()) + "/" + time.Now().Month().String() + "/" + strconv.Itoa(time.Now().Day())
		user.Subscription = timeNow
	}

	if user.Role == "" {
		user.Role = "user"
	}
	// 3. Save user in user's repo.
	if err := u.repo.AddUser(user); err != nil {
		return err
	}

	return nil
}

// GetByID returns a user according to given user id.
func (u *userSvc) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	return u.repo.GetUserByID(userID)
}

// GetUsers returns users registred in the app.
func (u *userSvc) GetUsers(ctx context.Context) (*entity.UsersResponse, error) {
	var offset, limit int64
	var department, filter string

	// Check if there is a valid page requested by the client.
	if !utils.IsValidPage(ctx.Value("offset").(string), ctx.Value("limit").(string)) {
		// If is not valid, then get a default page.
		offset, limit = utils.GetValidPage()
	} else {
		offset, _ = strconv.ParseInt(ctx.Value("offset").(string), 10, 64)
		limit, _ = strconv.ParseInt(ctx.Value("limit").(string), 10, 64)
	}

	department = ctx.Value("department").(string)
	filter = ctx.Value("filter").(string)

	totalItems, err := u.repo.GetAllUsersCount()
	if err != nil {
		return nil, err
	}
	items, err := u.repo.GetUsersByPage(offset, limit, department, filter)
	if err != nil {
		return nil, err
	}
	response := entity.UsersResponse{
		TotalItems: totalItems,
		Items:      items,
	}
	return &response, nil
}

// SignInWithPass pass the credentials needed to the auth service to sign in a user.
func (u *userSvc) SignInWithPass(c context.Context, creds *entity.StandardLoginCredentials) (*entity.AuthResponse, error) {
	// 1. SignIn user with creds
	idToken, err := u.authSvc.SignInWithPass(c, creds)
	if err != nil {
		logrus.Error("Step 1/6. Failed to SignInWithPass: " + err.Error())
		return nil, err
	}
	// 2. If there is no error, then get user from firestore
	user, err := u.repo.GetUserByEmail(creds.Email)
	if err != nil {
		// 2.1 If the user could make login in firebase but is not found in firestore.
		if errors.IsErrType(err) && err.(errors.Error).Kind.Code == errors.NotFound {
			// Then we need delete the account from firebase.
			u.authSvc.RemoveUser(idToken)
		}
		logrus.Error("Step 2/7. Failed to GetUserByEmail: " + err.Error())
		return nil, err
	}
	// 3. In order to persist only one session active
	// we remove any refresh token that the user may have
	if err := u.authSvc.RevokeUserTokens(user.ID); err != nil {
		logrus.Error("Step 3/7. Failed to revoke tokens: " + err.Error())
		return nil, err
	}

	// 4. Set claims(user info encrypted inside the token)
	claims := map[string]interface{}{
		"email":        user.Email,
		"subscription": user.Subscription,
		"fullName":     fmt.Sprintf("%s %s", user.Name, user.LastName),
		"birthday":     user.Birthday,
		"phone_number": user.PhoneNumber,
		"role":         user.Role,
	}
	// 5. Gen custom token with claims, info will be provided from the step 2
	// We need to set those claims for future request, we can read the JWT and get
	// User's information without requests to firestore
	token, err := u.authSvc.GenerateCustomToken(c, user.ID, claims)
	if err != nil {
		logrus.Error("Step 5/7. Failed to SignInWithClaims: " + err.Error())
		return nil, err
	}
	// 6. Once we have our custom token, in order to verify the JWT in each request
	// We need to exchange the custom token for a token id, which it will be readen
	// In the middleware to verify user's session.
	resp, err := u.authSvc.SignInWithTokenClaims(c, token)
	if err != nil {
		logrus.Error("Step 6/7. Failed to SignInWithClaims: " + err.Error())
		return nil, err
	}
	// 7. Returning response
	return &entity.AuthResponse{
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		User:         *user,
	}, nil
}
