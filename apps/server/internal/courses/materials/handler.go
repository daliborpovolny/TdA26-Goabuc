package materials

import (
	"fmt"
	"net/http"
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/handlers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	*handlers.Handler
	service    *Service
	staticPath string
}

func NewHandler(staticPath string, queries *db.Queries, isDeployed bool, service *Service) *Handler {
	return &Handler{
		handlers.NewHandler(queries, isDeployed),
		service,
		staticPath,
	}
}

func (h *Handler) ListMaterials(c echo.Context) error {
	r := h.NewReqCtx(c)
	courseId := c.Param("courseId")

	req := c.Request()

	mats, err := h.service.ListMaterials(courseId, req.Host, req.URL.Scheme, r.Ctx)
	if err != nil {
		if err == CourseNotFound {
			return r.Error(http.StatusBadRequest, "Unknown course id")
		}
		fmt.Println(err)
		return r.Error(http.StatusInternalServerError, "Failed to fetch materials from db")
	}

	return c.JSON(http.StatusOK, mats)
}

type CreateUrlMaterialRequest struct {
	MatType     string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func (h *Handler) collectFileMaterialParams(c echo.Context, r *handlers.RequestCtx) (*CreateFileMaterialParams, error) {

	courseId := c.Param("courseId")

	name := c.FormValue("name")
	if name == "" {
		return nil, r.Error(http.StatusBadRequest, "name of the material must be provided")
	}

	description := c.FormValue("description")

	fmt.Println("in file")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, r.Error(http.StatusBadRequest, "file is required")
	}

	req := c.Request()
	scheme := c.Scheme()
	host := req.Host

	return &CreateFileMaterialParams{
		fileHeader:  fileHeader,
		uuid:        uuid.NewString(),
		name:        name,
		description: description,
		courseId:    courseId,
		scheme:      scheme,
		host:        host,
	}, nil

}

func (h *Handler) handleCreateFileMaterialErrors(err error, r *handlers.RequestCtx) error {
	if err == TooBigMaterialFile {
		return r.Error(http.StatusBadRequest, "File is too big")
	}
	if err == FailedToOpenMaterialFile {
		return r.Error(http.StatusBadRequest, "Failed to open the file")
	}
	if err == ForbiddenFileType {
		return r.Error(http.StatusBadRequest, "Forbidden file type")
	}
	fmt.Println(err)
	return r.Error(http.StatusInternalServerError, "Creating file material failed")
}

func (h *Handler) CreateMaterial(c echo.Context) error {
	r := h.NewReqCtx(c)

	formMaterialType := c.FormValue("type")

	if formMaterialType == "file" {

		params, err := h.collectFileMaterialParams(c, r)
		if err != nil {
			return err
		}
		_, err = h.service.CreateFileMaterial(params, r.Ctx)
		if err != nil {
			return h.handleCreateFileMaterialErrors(err, r)
		}
		return r.JSONMsg(http.StatusCreated, "Created material")

	} else if formMaterialType != "" {
		return r.Error(http.StatusBadRequest, "Only file materials can be uploaded through a form")
	}

	var req CreateUrlMaterialRequest
	if err := c.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, "invalid create url material request")
	}

	if req.MatType == "url" {

		courseId := c.Param("courseId")

		if req.Name == "" {
			return r.Error(http.StatusBadRequest, "name of the material must be provided")
		}

		_, err := h.service.CreateUrlMaterial(CreateUrlMaterialParams{
			uuid:        uuid.NewString(),
			name:        req.Name,
			description: req.Description,
			courseId:    courseId,
			url:         "https://www.youtube.com/watch?v=3jL4S4X97sQ",
		}, r.Ctx)
		if err != nil {
			fmt.Println(err)
			return r.Error(http.StatusInternalServerError, "Failed to create url material")
		}
		return r.JSONMsg(http.StatusCreated, "Created url material")
	} else {
		return r.Error(http.StatusBadRequest, "Only url materials can be uploaded through json")
	}
}

func (h *Handler) UpdateMaterial(c echo.Context) error {
	r := h.NewReqCtx(c)

	materialId := c.Param("materialId")

	formMaterialType := c.FormValue("type")

	if formMaterialType == "file" {

		params, err := h.collectFileMaterialParams(c, r)
		if err != nil {
			return err
		}
		params.uuid = materialId // override new uuid created in collect

		_, err = h.service.UpdateFileMaterial(params, r.Ctx)
		if err != nil {
			return h.handleCreateFileMaterialErrors(err, r)
		}
		return r.JSONMsg(http.StatusCreated, "Updated material")

	} else if formMaterialType != "" {
		return r.Error(http.StatusBadRequest, "Only file materials can be uploaded through a form")
	}

	var req CreateUrlMaterialRequest
	if err := c.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, "invalid create url material request")
	}

	if req.MatType == "url" {

		courseId := c.Param("courseId")

		if req.Name == "" {
			return r.Error(http.StatusBadRequest, "name of the material must be provided")
		}

		_, err := h.service.CreateUrlMaterial(CreateUrlMaterialParams{
			uuid:        materialId,
			name:        req.Name,
			description: req.Description,
			courseId:    courseId,
			url:         "https://www.youtube.com/watch?v=3jL4S4X97sQ",
		}, r.Ctx)
		if err != nil {
			fmt.Println(err)
			return r.Error(http.StatusInternalServerError, "Failed to create url material")
		}
		return r.JSONMsg(http.StatusCreated, "Created url material")
	} else {
		return r.Error(http.StatusBadRequest, "Only url materials can be uploaded through json")
	}
}

func (h *Handler) DeleteMaterial(c echo.Context) error {
	r := h.NewReqCtx(c)

	// courseId := c.Param("courseId")
	materialId := c.Param("materialId")

	err := h.service.DeleteMaterial(materialId, r.Ctx)
	if err != nil {
		if err == CourseNotFound {
			return r.Error(http.StatusBadRequest, "Material not found")
		}
		return r.Error(http.StatusInternalServerError, "Failed to delete material")
	}

	return r.JSONMsg(http.StatusNoContent, "Material deleted")
}
