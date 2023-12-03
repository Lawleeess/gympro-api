package user

import (
	"context"
	"fmt"
	"math"
	"mime/multipart"
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
	user.Url = "https://firebasestorage.googleapis.com/v0/b/gympro-400622.appspot.com/o/users%2Fuser_default.png?alt=media&token=f5434b37-9f27-4f1d-9b00-0cbf69f24c2e"
	if user.UserRole == "" {
		user.UserRole = "user"
	} else if user.UserRole == "admin" {
		moduleUserManagement := entity.Module{
			"name": "userManagement",
			"role": "admin",
		}
		moduleRoutinesManagement := entity.Module{
			"name": "routinesManagement",
			"role": "admin",
		}

		user.Modules = append(user.Modules, moduleUserManagement)
		user.Modules = append(user.Modules, moduleRoutinesManagement)
	}

	if user.Subscription == "" {
		timeNow := time.Now()
		yesterday := timeNow.AddDate(0, 0, -1).Format("2006-01-02")

		// timeNow := strconv.Itoa(time.Now().Year()) + "/" + time.Now().Month().String() + "/" + strconv.Itoa(time.Now().Day())
		user.Subscription = yesterday
	}

	// 3. Save user in user's repo.
	if err := u.repo.AddUser(user); err != nil {
		return err
	}

	return nil
}

// GetByID returns a user according to given user id.
func (u *userSvc) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	return u.repo.GetUserByID(userID)
}

// GetUsers returns users registred in the app.
func (u *userSvc) GetUsers(ctx context.Context) (*entity.UsersResponse, error) {
	var offset, limit int64

	// Check if there is a valid page requested by the client.
	if !utils.IsValidPage(ctx.Value("offset").(string), ctx.Value("limit").(string)) {
		// If is not valid, then get a default page.
		offset, limit = utils.GetValidPage()
	} else {
		offset, _ = strconv.ParseInt(ctx.Value("offset").(string), 10, 64)
		limit, _ = strconv.ParseInt(ctx.Value("limit").(string), 10, 64)
	}

	userRole := ctx.Value("user_role").(string)
	filter := ctx.Value("filter").(string)

	totalItems, err := u.repo.GetAllUsersCount()
	if err != nil {
		return nil, err
	}
	items, err := u.repo.GetUsersByPage(offset, limit, userRole, filter)
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
		"user_id":               user.ID,
		"email":                 user.Email,
		"subscription":          user.Subscription,
		"fullName":              fmt.Sprintf("%s %s", user.Name, user.LastName),
		"birthday":              user.Birthday,
		"phone_number":          user.PhoneNumber,
		"modulesWithPermission": user.Modules,
		"url_image":             user.Url,
		"user_role":             user.UserRole,
		"userProgress":          user.UserProgress,
		"userGoals":             user.UserGoals,
		"userRoutine":           user.UserRoutine,
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

	curretDayTime := time.Now()
	curretDay := curretDayTime.Format("2006-01-02")
	dateCurrent, _ := time.Parse("2006-01-02", curretDay)
	dateSubs, _ := time.Parse("2006-01-02", user.Subscription)

	if dateSubs.Before(dateCurrent) && user.UserRole != "admin" {
		x := make([]map[string]interface{}, 0)
		user.Modules = x
		errUpdate := u.repo.UpdateUser(user.ID, user)
		if errUpdate != nil {
			return nil, errUpdate
		}

	} else if user.Modules == nil {
		u.UpdateUser(user.ID, user)
	}

	// 7. Returning response
	return &entity.AuthResponse{
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		User:         *user,
	}, nil
}

func (u *userSvc) UpdateUser(userID string, user *entity.User) error {

	curretDayTime := time.Now()
	curretDay := curretDayTime.Format("2006-01-02")
	dateCurrent, _ := time.Parse("2006-01-02", curretDay)
	dateSubs, _ := time.Parse("2006-01-02", user.Subscription)

	if dateSubs.After(dateCurrent) && user.Modules == nil {
		moduleRoutines := entity.Module{
			"name": "routinesCalendar",
			"role": "viewer",
		}
		modulePersonalGoals := entity.Module{
			"name": "personalGoals",
			"role": "viewer",
		}

		user.Modules = append(user.Modules, moduleRoutines)
		user.Modules = append(user.Modules, modulePersonalGoals)
	}

	errUpdate := u.repo.UpdateUser(userID, user)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (u *userSvc) UpdateImageUser(img multipart.File, userID string) error {

	urlImg, err := u.authSvc.UpdateUserImage(img, userID)
	if err != nil {
		return err
	}

	errUpdate := u.repo.UpdateImageUser(userID, urlImg)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

// DeleteUser removes the given user from the app.
func (u *userSvc) DeleteUser(ctx context.Context, userID string) error {
	currentUser := ctx.Value("userID").(string)

	// Check if the user making the request is trying to remove himself.
	if currentUser == userID {
		return errors.Build(
			errors.Operation("userService.DeleteUser"),
			errors.Forbidden,
			errors.Message("You can't remove your own user."),
		)
	}
	return u.repo.DeleteUser(userID)
}

// VerifyToken virifies if the given token is valid and return its claims.
func (u *userSvc) VerifyToken(token string) (map[string]interface{}, error) {
	scope := errors.Operation("authService.VerifyToken")

	authenticatedToken, err := u.authSvc.VerifyToken(token)
	if err != nil {
		return nil, errors.Build(
			scope,
			errors.Unauthorized,
			errors.Message("Invalid session"),
		)
	}
	return authenticatedToken.Claims, nil
}

func (u *userSvc) SaveUserProgress(userID string, userProgress *entity.UserProgress) (*entity.UserGoals, error) {
	scope := errors.Operation("userServuce.SaveUserProgress")

	err := u.repo.SaveUserProgress(userID, userProgress)
	if err != nil {
		return nil, errors.Build(
			scope,
			errors.Unauthorized,
			errors.Message("Error saving user progress"),
		)
	}

	userGoals := &entity.UserGoals{}

	userGoals.IMC = fmt.Sprintf("%.2f", userProgress.Weight/(math.Pow((float64(userProgress.Height)*.01), 2)))

	if userProgress.Gender == "hombre" {
		userGoals.BMR = fmt.Sprintf("%.2f", ((10 * userProgress.Weight) + (float64(6.25) * float64(userProgress.Height)) - float64(5*userProgress.Age) + 5))
	} else {
		userGoals.BMR = fmt.Sprintf("%.2f", ((10 * userProgress.Weight) + (float64(6.25) * float64(userProgress.Height)) - float64(5*userProgress.Age) - 161))
	}

	var bmr, activity, tdee, goal, protein, carbohydrates, fat, userGoal float64

	if bmr, err = strconv.ParseFloat(userGoals.BMR, 64); err != nil {
		return nil, err
	}

	if activity, err = strconv.ParseFloat(userProgress.Activity, 64); err != nil {
		return nil, err
	}

	userGoals.TDEE = fmt.Sprintf("%.2f", bmr*activity)

	if tdee, err = strconv.ParseFloat(userGoals.TDEE, 64); err != nil {
		return nil, err
	}

	if goal, err = strconv.ParseFloat(userProgress.Goal, 64); err != nil {
		return nil, err
	}

	userGoals.Goal = fmt.Sprintf("%.2f", tdee*goal)

	if userGoal, err = strconv.ParseFloat(userGoals.Goal, 64); err != nil {
		return nil, err
	}

	if userProgress.Goal == "1" { //mantener el peso
		protein = (userGoal * 0.25) / 4
		carbohydrates = (userGoal * 0.5) / 4
		fat = (userGoal * 0.25) / 9
	} else if userProgress.Goal == "1.21" { //aumentar peso
		protein = (userGoal * 0.35) / 4
		carbohydrates = (userGoal * 0.45) / 4
		fat = (userGoal * 0.2) / 9
	} else if userProgress.Goal == "0.79" { //bajar peso
		protein = (userGoal * 0.45) / 4
		carbohydrates = (userGoal * 0.35) / 4
		fat = (userGoal * 0.2) / 9
	} else if userProgress.Goal == "1.10" { //aumentar peso leve
		protein = (userGoal * 0.30) / 4
		carbohydrates = (userGoal * 0.45) / 4
		fat = (userGoal * 0.25) / 9
	} else if userProgress.Goal == "0.9" { //bajar peso leve
		protein = (userGoal * 0.40) / 4
		carbohydrates = (userGoal * 0.35) / 4
		fat = (userGoal * 0.25) / 9
	}

	userGoals.Protein = fmt.Sprintf("%.2f", protein)
	userGoals.Carbs = fmt.Sprintf("%.2f", carbohydrates)
	userGoals.Fat = fmt.Sprintf("%.2f", fat)

	err2 := u.repo.SaveUserGoals(userID, userGoals)
	if err2 != nil {
		return nil, errors.Build(
			scope,
			errors.Unauthorized,
			errors.Message("Error saving user goals"),
		)
	}

	return userGoals, nil
}

func (u *userSvc) AddRoutineToUser(userID string, userRoutine *entity.UserRoutine) error {
	return u.repo.AddRoutineToUser(userID, userRoutine)
}

func (u userSvc) VerifyOrRecoverEmail(ctx context.Context, creds *entity.UserRequestType) (string, error) {
	return u.authSvc.VerifyOrRecoverEmail(ctx, creds)
}

func (u userSvc) VerifyOobCode(ctx context.Context, creds *entity.OobCode) (bool, error) {
	return u.authSvc.VerifyOobCode(ctx, creds)
}
