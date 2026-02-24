package courses

import (
	"net/http"
	"time"

	db "tourbackend/internal/database/gen"
	"tourbackend/internal/handlers"
	"tourbackend/internal/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CourseHandler struct {
	*handlers.Handler
	service *Service
}

func NewCourseHandler(queries *db.Queries, isDeployed bool, service *Service) *CourseHandler {
	return &CourseHandler{
		handlers.NewHandler(queries, isDeployed),
		service,
	}
}

type CreateCourseRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateCourseResponse struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func (h *CourseHandler) CreateCourse(c echo.Context) error {
	r := h.NewReqCtx(c)

	var req CreateCourseRequest
	if err := c.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, "invalid request")
	}

	unixTime := time.Now().Unix()

	params := db.CreateCourseParams{
		Uuid:        uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   unixTime,
		UpdatedAt:   unixTime,
	}

	course, err := h.service.CreateCourse(params, r.Ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, CreateCourseResponse{
		Uuid:        course.Uuid,
		Name:        course.Name,
		Description: course.Description,
		CreatedAt:   utils.UnixToIso(course.CreatedAt),
		UpdatedAt:   utils.UnixToIso(course.UpdatedAt),
	})
}

func (h *CourseHandler) GetCourse(c echo.Context) error {
	r := h.NewReqCtx(c)

	courseId := c.Param("courseId")

	if courseId == "" {
		return r.Error(http.StatusBadRequest, "Must specify the course uuid.")
	}

	req := c.Request()
	courseDetail, err := h.service.GetCourse(courseId, req.Host, c.Scheme(), r.Ctx)
	if err != nil {
		if err == ErrFailedToFetchCourse {
			return r.Error(http.StatusInternalServerError, "Failed to fetch from the database")
		}
		if err == ErrCourseNotFound {
			return r.Error(http.StatusNotFound, "Unknown courseId")
		}
		return r.Error(http.StatusInternalServerError, "Failed to get course")
	}

	return c.JSON(http.StatusOK, *courseDetail)

}

type UpdateCourseRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCourseResponse struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func (h *CourseHandler) UpdateCourse(c echo.Context) error {
	r := h.NewReqCtx(c)

	courseId := c.Param("courseId")
	if courseId == "" {
		return r.Error(http.StatusBadRequest, "invalid request, missing course id")
	}

	var req UpdateCourseRequest
	if err := c.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, "invalid request")
	}

	unixTime := time.Now().Unix()

	updateParams := db.UpdateCourseParams{
		Name:        req.Name,
		Description: req.Description,
		UpdatedAt:   unixTime,
		Uuid:        courseId,
	}

	course, err := h.service.UpdateCourse(updateParams, r.Ctx)
	if err != nil {
		if err == ErrCourseNotFound {
			return r.Error(http.StatusNotFound, "The requested resource was not found.")
		}
		return r.Error(http.StatusInternalServerError, "Failed to update the course.")
	}

	return c.JSON(http.StatusCreated, UpdateCourseResponse{
		Uuid:        courseId,
		Name:        course.Name,
		Description: course.Description,
		CreatedAt:   utils.UnixToIso(course.CreatedAt),
		UpdatedAt:   utils.UnixToIso(course.UpdatedAt),
	})
}

func (h *CourseHandler) DeleteCourse(c echo.Context) error {
	r := h.NewReqCtx(c)

	courseId := c.Param("courseId")
	if courseId == "" {
		return r.Error(http.StatusBadRequest, "courseId must be provided as path parameter")
	}

	err := h.service.DeleteCourse(courseId, r.Ctx)
	if err != nil {
		if err == ErrFailedToFetchCourse {
			return r.Error(http.StatusNotFound, "The requested resource was not found.")
		}
		return r.Error(http.StatusInternalServerError, "Failed to delete the course.")
	}

	return c.NoContent(http.StatusNoContent)
}

type ListAllCoursesResponse struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func (h *CourseHandler) ListAllCourses(c echo.Context) error {
	r := h.NewReqCtx(c)

	courses, err := h.service.ListAllCourses(r.Ctx)
	if err != nil {
		r.ServerError(err)
	}

	formattedCourses := make([]ListAllCoursesResponse, len(courses))
	for i, course := range courses {
		formattedCourses[i] = ListAllCoursesResponse{
			Uuid:        course.Uuid,
			Name:        course.Name,
			Description: course.Description,
			CreatedAt:   utils.UnixToIso(course.CreatedAt),
			UpdatedAt:   utils.UnixToIso(course.UpdatedAt),
		}
	}

	return c.JSON(http.StatusOK, formattedCourses)
}

type ChangeCourseStateRequest struct {
	State string `json:"state"`
}

func (h *CourseHandler) ChangeCourseState(c echo.Context) error {
	r := h.NewReqCtx(c)

	var req ChangeCourseStateRequest
	if err := c.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, "invalid request")
	}

	courseId := r.Echo.Param("courseId")

	_, err := h.service.ChangeCourseState(courseId, req.State, r.Ctx)
	if err != nil {
		if err == ErrBadCourseState {
			return r.Error(http.StatusBadRequest, "Invalid course state")
		}
		return r.ServerError(err)
	}

	return r.JSONMsg(http.StatusCreated, "New course state set")
}

// modules

type ChangeModuleStateRequest struct {
	State string `json:"state"`
}

func (h *CourseHandler) ChangeModuleState(c echo.Context) error {
	r := h.NewReqCtx(c)

	var req ChangeModuleStateRequest
	if err := c.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, "invalid request, must provide state")
	}

	courseId := r.Echo.Param("courseId")
	moduleId := r.Echo.Param("moduleId")

	_, err := h.service.ChangeModuleState(courseId, moduleId, req.State, r.Ctx)
	if err != nil {
		if err == ErrBadModuleState {
			return r.Error(http.StatusBadRequest, "Invalid module state")
		}
		return r.ServerError(err)
	}

	return r.JSONMsg(http.StatusCreated, "New module state set")
}

type CreateNewModuleRequest struct {
	Name string `json:"name"`
}

func (h *CourseHandler) CreateNewModule(c echo.Context) error {
	r := h.NewReqCtx(c)

	var req CreateNewModuleRequest
	if err := c.Bind(&req); err != nil {
		return r.Error(http.StatusBadRequest, "invalid request, must provide name")
	}

	courseId := r.Echo.Param("courseId")
	moduleId := uuid.NewString()

	_, err := h.service.CreateModule(courseId, moduleId, req.Name, r.Ctx)
	if err != nil {
		return r.ServerError(err)
	}

	return r.JSONMsg(http.StatusCreated, "module created")
}

func (h *CourseHandler) GetModule(c echo.Context) error {
	r := h.NewReqCtx(c)

	courseId := r.Echo.Param("courseId")
	moduleId := r.Echo.Param("moduleId")

	module, err := h.service.GetModule(courseId, moduleId, r.Ctx)
	if err != nil {
		return r.ServerError(err)
	}

	return c.JSON(http.StatusOK, module)
}
