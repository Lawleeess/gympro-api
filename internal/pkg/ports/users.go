package ports

import (
	"context"
	"mime/multipart"

	"github.com/CValier/gympro-api/internal/pkg/entity"
)

// UsersRepository is the signature to perform operations related to user's CRUD operations
type UsersRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUsersByPage(offset, limit int64, userRole, filter string) ([]*entity.User, error)
	GetUsersByPageActive(offset, limit int64, userRole, filter string) ([]*entity.User, error)
	GetUserByID(userID string) (*entity.User, error)
	AddUser(user *entity.User) error
	DeleteUser(userID string) error
	GetAllUsersCount() (int, error)
	UpdateUser(userID string, user *entity.User) error
	SaveUserProgress(userID string, userProgress *entity.UserProgress) error
	SaveUserGoals(userID string, userGoals *entity.UserGoals) error
	UpdateImageUser(userID string, url string) error

	//Routine
	AddRoutine(routine *entity.Routine) error
	UpdateImageRoutine(id string, url string) error
	AddRoutineToUser(userID string, userRoutine *entity.UserRoutine) error
	GetRoutines(muscle_group string) ([]entity.Routine, error)
}

// UserService is the signature to perform business logic over the user resource.
type UserService interface {
	SignInWithPass(ctx context.Context, credentials *entity.StandardLoginCredentials) (*entity.AuthResponse, error)
	VerifyOrRecoverEmail(ctx context.Context, creds *entity.UserRequestType) (string, error)
	VerifyOobCode(ctx context.Context, creds *entity.UserVerify) error

	GetUsersActive(ctx context.Context) (*entity.UsersResponse, error)
	GetUsers(ctx context.Context) (*entity.UsersResponse, error)
	GetUserByID(ctx context.Context, userID string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(userID string, user *entity.User) error
	DeleteUser(ctx context.Context, userID string) error
	SaveUserProgress(userID string, userProgress *entity.UserProgress) (*entity.UserGoals, error)
	UpdateImageUser(img multipart.File, userID string) error
	VerifyToken(token string) (map[string]interface{}, error)

	//Routine
	AddRoutineToUser(userID string, userRoutine *entity.UserRoutine) error
}
