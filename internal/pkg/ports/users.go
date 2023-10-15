package ports

import (
	"context"
	"mime/multipart"

	"github.com/CValier/gympro-api/internal/pkg/entity"
)

// UsersRepository is the signature to perform operations related to user's CRUD operations
type UsersRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUsersByPage(offset, limit int64, department, filter string) ([]*entity.User, error)
	GetUserByID(userID string) (*entity.User, error)
	AddUser(user *entity.User) error
	GetAllUsersCount() (int, error)
	UpdateUser(userID string, user *entity.User) error
	UpdateImageUser(userID string, url string) error
}

// UserService is the signature to perform business logic over the user resource.
type UserService interface {
	SignInWithPass(ctx context.Context, credentials *entity.StandardLoginCredentials) (*entity.AuthResponse, error)
	GetUsers(ctx context.Context) (*entity.UsersResponse, error)
	GetByID(ctx context.Context, userID string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(userID string, user *entity.User) error
	UpdateImageUser(img multipart.File, userID string) error
	VerifyToken(token string) (map[string]interface{}, error)
}
