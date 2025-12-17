package materials

import (
	"errors"
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

var errNoName = errors.New("material has no name")
var errNoFile = errors.New("no file")

func (h *Handler) collectFileMaterialParams(c echo.Context) (*CreateFileMaterialParams, error) {

	courseId := c.Param("courseId")

	name := c.FormValue("name")
	if name == "" {
		return nil, errNoName
	}

	description := c.FormValue("description")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, errNoFile
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

		params, err := h.collectFileMaterialParams(c)
		if err != nil {
			return err
		}
		dbMat, err := h.service.CreateFileMaterial(params, r.Ctx)
		if err != nil {
			return h.handleCreateFileMaterialErrors(err, r)
		}
		return c.JSON(http.StatusCreated, FileMaterial{
			Uuid:        dbMat.Uuid,
			Type:        "file",
			Name:        dbMat.Name,
			Description: dbMat.Description,
			FileUrl:     dbMat.Url,
		})

	} else if formMaterialType != "" {
		return r.Error(http.StatusBadRequest, "Only file materials can be uploaded through a form")
	}

	var req CreateUrlMaterialRequest
	if err := c.Bind(&req); err != nil {
		fmt.Println(err)
		return r.Error(http.StatusBadRequest, "invalid create url material request")
	}

	if req.MatType == "url" {

		courseId := c.Param("courseId")

		if req.Name == "" {
			return r.Error(http.StatusBadRequest, "name of the material must be provided")
		}

		dbMat, err := h.service.CreateUrlMaterial(CreateUrlMaterialParams{
			uuid:        uuid.NewString(),
			name:        req.Name,
			description: req.Description,
			courseId:    courseId,
			url:         req.Url,
		}, r.Ctx)
		if err != nil {
			fmt.Println(err)
			return r.Error(http.StatusInternalServerError, "Failed to create url material")
		}

		return c.JSON(http.StatusOK, UrlMaterial{
			Uuid:        dbMat.Uuid,
			Type:        "url",
			Name:        dbMat.Name,
			Description: dbMat.Description,
			Url:         dbMat.Url,
			FaviconUrl:  h.service.deriveFaviconUrl(dbMat.Url),
		})

	} else {
		return r.Error(http.StatusBadRequest, "Only url materials can be uploaded through json")
	}
}

func (h *Handler) UpdateMaterial(c echo.Context) error {
	r := h.NewReqCtx(c)

	materialId := c.Param("materialId")

	params, err := h.collectFileMaterialParams(c)
	fmt.Println(params, err, "par and err")
	if err == nil {
		fmt.Println("nill err", err)
		if err == errNoFile {
			return r.Error(http.StatusBadRequest, "file must be included")
		}

		params.uuid = materialId // override new uuid created in collect

		dbMat, err := h.service.UpdateFileMaterial(params, r.Ctx)
		if err != nil {
			return h.handleCreateFileMaterialErrors(err, r)
		}

		return c.JSON(http.StatusOK, FileMaterial{
			Uuid:        dbMat.Uuid,
			Type:        "file",
			Name:        dbMat.Name,
			Description: dbMat.Description,
			FileUrl:     dbMat.Url,
		})
	}

	fmt.Println("attempting json")

	var req CreateUrlMaterialRequest
	if err = c.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, "invalid update material request")
	}

	courseId := c.Param("courseId")

	if req.Name == "" {
		return r.Error(http.StatusBadRequest, "name of the material must be provided")
	}

	dbMat, err := h.service.UpdateUrlMaterial(&CreateUrlMaterialParams{
		uuid:        materialId,
		name:        req.Name,
		description: req.Description,
		courseId:    courseId,
		url:         req.Url,
	}, r.Ctx)
	if err != nil {
		fmt.Println(err)
		return r.Error(http.StatusInternalServerError, "Failed to update url material")
	}

	return c.JSON(http.StatusOK, UrlMaterial{
		Uuid:        dbMat.Uuid,
		Type:        "url",
		Name:        dbMat.Name,
		Description: dbMat.Description,
		Url:         dbMat.Url,
		FaviconUrl:  h.service.deriveFaviconUrl(dbMat.Url),
	})
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

	return c.NoContent(http.StatusNoContent)
}
