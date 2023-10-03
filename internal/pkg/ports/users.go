package ports

import (
	"context"

	"github.com/CValier/gympro-api/internal/pkg/entity"
)

// UsersRepository is the signature to perform operations related to user's CRUD operations
type UsersRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUsersByPage(offset, limit int64, department, filter string) ([]*entity.User, error)
	GetUserByID(userID string) (*entity.User, error)
	AddUser(user *entity.User) error
	GetAllUsersCount() (int, error)
}

// UserService is the signature to perform business logic over the user resource.
type UserService interface {
	SignInWithPass(ctx context.Context, credentials *entity.StandardLoginCredentials) (*entity.AuthResponse, error)
	GetUsers(ctx context.Context) (*entity.UsersResponse, error)
	GetByID(ctx context.Context, userID string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
}
