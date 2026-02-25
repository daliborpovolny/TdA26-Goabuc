package materials

import (
	"net/http"
	"strconv"
	"strings"
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/handlers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	*handlers.Handler
	service      *Service
	pathToStatic string
}

func NewHandler(pathToStatic string, service *Service, queries *db.Queries, isDeployed bool) *Handler {
	return &Handler{
		handlers.NewHandler(queries, isDeployed),
		service,
		pathToStatic,
	}
}

// List Materials

func (h *Handler) ListMaterials(c echo.Context) error {
	r := h.NewReqCtx(c)

	courseId := c.Param("courseId")
	req := c.Request()

	mats, err := h.service.ListMaterials(courseId, req.Host, req.URL.Scheme, r.Ctx)
	if err != nil {
		if err == ErrCourseNotFound {
			return r.Error(http.StatusNotFound, "Unknown course id")
		}
		r.ServerError(err)
	}

	return c.JSON(http.StatusOK, mats)
}

// Create Material

func (h *Handler) CreateMaterial(c echo.Context) error {
	r := h.NewReqCtx(c)

	contentType := c.Request().Header["Content-Type"]

	if strings.Contains(contentType[0], "multipart/form-data") {
		return h.createFileMaterial(r)

	} else if strings.Contains(contentType[0], "application/json") {
		return h.createUrlMaterial(r)
	} else {
		return r.Error(http.StatusBadRequest, "bad request")
	}
}

type CreateFileMaterialRequest struct {
	CourseId string `param:"courseId"`

	MatType string `form:"type"`
	Name    string `form:"name"`

	Description string `form:"description"`

	ModuleId    *string `form:"moduleId"`
	ModuleOrder *int    `form:"moduleOrder"`
}

func (h *Handler) createFileMaterial(r *handlers.RequestCtx) error {
	var req CreateFileMaterialRequest
	if err := r.Echo.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, err.Error())
	}

	if req.MatType != "file" {
		return r.Error(http.StatusBadRequest, "only file material can be created through form")
	}

	if req.Name == "" {
		return r.Error(http.StatusBadRequest, "name is required")
	}

	file, err := r.Echo.FormFile("file")
	if err != nil {
		return r.Error(http.StatusBadRequest, "file is required")
	}

	httpReq := r.Echo.Request()

	mat, err := h.service.CreateFileMaterial(&req, uuid.NewString(), file, r.Echo.Scheme(), httpReq.Host, r.Ctx)
	if err != nil {
		if err == ErrFileTooBig {
			return r.Error(http.StatusBadRequest, "file is too big")
		}
		if err == ErrFileTypeForbidden {
			return r.Error(http.StatusBadRequest, "file type forbidden")
		}
		return r.ServerError(err)
	}

	if req.ModuleId != nil {
		if req.ModuleOrder == nil {
			return r.Error(http.StatusBadRequest, "module order must be provided along with moduleId")
		}

		moduleId := *req.ModuleId
		order := *req.ModuleOrder

		_, err := h.service.AssignMaterialToModule(mat.GetUuid(), moduleId, order, r.Ctx)
		if err != nil {
			return r.ServerError(err)
		}
	} else if !MATERIAL_CAN_EXIST_ALONE {
		return r.Error(http.StatusBadRequest, "material must always be part of a module")
	}

	return r.Echo.JSON(http.StatusCreated, mat)
}

type CreateUrlMaterialRequest struct {
	CourseId string `param:"courseId"`

	MatType string `json:"type"`
	Name    string `json:"name"`
	Url     string `json:"url"`

	Description string `json:"description"`

	ModuleId    *string `json:"moduleId"`
	ModuleOrder *int    `json:"moduleOrder"`
}

func (h *Handler) createUrlMaterial(r *handlers.RequestCtx) error {
	var req CreateUrlMaterialRequest
	if err := r.Echo.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, err.Error())
	}

	if req.MatType != "url" {
		return r.Error(http.StatusBadRequest, "only url material can be created through form")
	}

	if req.Name == "" {
		return r.Error(http.StatusBadRequest, "name is required")
	}

	mat, err := h.service.CreateUrlMaterial(req, uuid.NewString(), r.Ctx)
	if err != nil {
		return r.ServerError(err)
	}

	if req.ModuleId != nil {
		if req.ModuleOrder == nil {
			return r.Error(http.StatusBadRequest, "module order must be provided along with moduleId")
		}

		moduleId := *req.ModuleId
		order := *req.ModuleOrder

		_, err := h.service.AssignMaterialToModule(mat.GetUuid(), moduleId, order, r.Ctx)
		if err != nil {
			return r.ServerError(err)
		}
	} else if !MATERIAL_CAN_EXIST_ALONE {
		return r.Error(http.StatusBadRequest, "material must always be part of a module")
	}

	return r.Echo.JSON(http.StatusCreated, mat)
}

// Update Materials

func (h *Handler) UpdateMaterial(c echo.Context) error {
	r := h.NewReqCtx(c)

	contentType := c.Request().Header["Content-Type"]

	if strings.Contains(contentType[0], "multipart/form-data") {
		return h.updateFileMaterial(r)

	} else if strings.Contains(contentType[0], "application/json") {
		return h.updateUrlMaterial(r)
	} else {
		return r.Error(http.StatusBadRequest, "bad request")
	}
}

type UpdateFileMaterialRequest struct {
	CourseId   string `param:"courseId"`
	MaterialId string `param:"materialId"`

	Name        *string `form:"name"`
	Description *string `form:"description"`
}

func (h *Handler) updateFileMaterial(r *handlers.RequestCtx) error {
	var req UpdateFileMaterialRequest
	if err := r.Echo.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, err.Error())
	}

	file, err := r.Echo.FormFile("file")
	if err != nil {
		file = nil
	}

	httpReq := r.Echo.Request()

	mat, err := h.service.UpdateFileMaterial(&req, file, r.Echo.Scheme(), httpReq.Host, r.Ctx)
	if err != nil {
		if err == ErrFileTooBig {
			return r.Error(http.StatusBadRequest, "file is too big")
		}
		if err == ErrFileTypeForbidden {
			return r.Error(http.StatusBadRequest, "file type forbidden")
		}
		return r.ServerError(err)
	}

	return r.Echo.JSON(http.StatusCreated, mat)
}

type UpdateUrlMaterialRequest struct {
	CourseId   string `param:"courseId"`
	MaterialId string `param:"materialId"`

	Name        *string `json:"name"`
	Url         *string `json:"url"`
	Description *string `json:"description"`
}

func (h *Handler) updateUrlMaterial(r *handlers.RequestCtx) error {
	var req UpdateUrlMaterialRequest
	if err := r.Echo.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, err.Error())
	}

	mat, err := h.service.UpdateUrlMaterial(&req, r.Ctx)
	if err != nil {
		return r.ServerError(err)
	}

	return r.Echo.JSON(http.StatusCreated, mat)
}

func (h *Handler) DeleteMaterial(c echo.Context) error {
	r := h.NewReqCtx(c)

	materialId := c.Param("materialId")

	err := h.service.DeleteMaterial(materialId, r.Ctx)
	if err != nil {
		if err == ErrCourseNotFound {
			return r.Error(http.StatusBadRequest, "Material not found")
		}
		return r.ServerError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ChangeMaterialInModuleOrder(c echo.Context) error {
	r := h.NewReqCtx(c)

	materialId := c.Param("materialId")
	moduleId := c.Param("moduleId")

	orderStr := c.Param("order")
	order, err := strconv.Atoi(orderStr)
	if err != nil {
		return r.Error(http.StatusBadRequest, "order must be a number")
	}

	_, err = h.service.ChangeMaterialInModuleOrder(materialId, moduleId, order, r.Ctx)
	if err != nil {
		return r.ServerError(err)
	}
	return r.JSONMsg(http.StatusCreated, "changed the order")
}
