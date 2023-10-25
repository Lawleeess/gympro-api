package routines

import (
	"net/http"

	"github.com/CValier/gympro-api/internal/pkg/entity"
	"github.com/CValier/gympro-api/internal/pkg/ports"
	"github.com/epa-datos/errors"
	"github.com/gin-gonic/gin"
)

type routinesHandler struct {
	routineService ports.RoutineService
}

func newHandler(service ports.RoutineService) *routinesHandler {
	return &routinesHandler{
		routineService: service,
	}
}

func (u *routinesHandler) addRoutine(c *gin.Context) {
	routine := &entity.Routine{}

	if err := c.Bind(&routine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid format: " + err.Error(),
		})
		return
	}

	err := u.routineService.AddRoutine(routine)
	if err != nil {
		errors.JSON(c, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (u *routinesHandler) updateRoutineUser(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.Build(
			errors.Message("Failed to get image: "+err.Error()),
		))
		return
	}

	errUpdate := u.routineService.UpdateRoutineImage(file, c.Param("id"))
	if errUpdate != nil {
		errors.JSON(c, errUpdate)
		return
	}

	c.JSON(http.StatusOK, nil)
}
