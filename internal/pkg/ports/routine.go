package ports

import (
	"mime/multipart"

	"github.com/CValier/gympro-api/internal/pkg/entity"
)

type RoutineService interface {
	AddRoutine(routine *entity.Routine) error
	UpdateRoutineImage(img multipart.File, id string) error
	GetRoutines(muscle_group string) ([]entity.Routine, error)
}
