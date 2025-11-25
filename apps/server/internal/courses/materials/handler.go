package materials

import (
	"net/http"
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/handlers"

	"github.com/labstack/echo/v4"
)

type MaterialsHandlers struct {
	*handlers.Handler
}

func NewMaterialsHandlers(queries *db.Queries, isDeployed bool) *MaterialsHandlers {
	return &MaterialsHandlers{
		handlers.NewHandler(queries, isDeployed),
	}
}

func (h *MaterialsHandlers) ListMaterials(c echo.Context) error {
	courseId := c.Param("courseId")

	return c.String(http.StatusOK, "List materials, Coures: "+courseId)
}

func (h *MaterialsHandlers) CreateMaterial(c echo.Context) error {
	courseId := c.Param("courseId")

	return c.String(http.StatusOK, "Create material, Coures: "+courseId)
}

func (h *MaterialsHandlers) UpdateMaterial(c echo.Context) error {
	courseId := c.Param("courseId")
	materialId := c.Param("materialId")

	return c.String(http.StatusOK, "Update material, Course: "+courseId+" Material:"+materialId)
}

func (h *MaterialsHandlers) DeleteMaterial(c echo.Context) error {
	courseId := c.Param("courseId")
	materialId := c.Param("materialId")

	return c.String(http.StatusOK, "Delete material: Course:"+courseId+" Material:"+materialId)
}
