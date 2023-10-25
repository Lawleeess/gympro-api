package routine

import (
	"mime/multipart"

	"github.com/CValier/gympro-api/internal/pkg/entity"
	"github.com/CValier/gympro-api/internal/pkg/ports"
)

type routineSvc struct {
	repo    ports.UsersRepository
	authSvc ports.AuthProvider
}

// NewUserService returns an instance for user service
func NewUserService(repo ports.UsersRepository, firebAuth ports.AuthProvider) *routineSvc {
	return &routineSvc{
		repo:    repo,
		authSvc: firebAuth,
	}
}

func (u *routineSvc) AddRoutine(routine *entity.Routine) error {
	errUpdate := u.repo.AddRoutine(routine)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (u *routineSvc) UpdateRoutineImage(img multipart.File, id string) error {

	urlImg, err := u.authSvc.UpdateRoutineImage(img, id)
	if err != nil {
		return err
	}

	errUpdate := u.repo.UpdateImageRoutine(id, urlImg)
	if errUpdate != nil {
		return errUpdate
	}

	return nil
}
