package quizzes

import (
	"net/http"
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/handlers"

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

func (h *Handler) ListQuizzes(c echo.Context) error {
	return nil
}

func (h *Handler) CreateQuizz(c echo.Context) error {
	r := h.NewReqCtx(c)

	courseId := r.Echo.Param("courseId")

	var quiz Quiz
	if err := c.Bind(&quiz); err != nil {
		return r.Error(http.StatusBadRequest, "bad request")
	}

	ok := h.service.validateQuestions(quiz.Questions)
	if !ok {
		return r.Error(http.StatusBadRequest, "invalid question")
	}

	dbQuiz, err := h.service.CreateQuiz(quiz, courseId, r.Ctx)
	if err != nil {
		return r.ServerError(err)
	}
	return r.Echo.JSON(http.StatusCreated, dbQuiz)
}

func (h *Handler) GetQuizz(c echo.Context) error {
	return nil
}

func (h *Handler) UpdateQuizz(c echo.Context) error {
	return nil
}

func (h *Handler) DeleteQuizz(c echo.Context) error {
	return nil
}

func (h *Handler) SubmitQuizAnswers(c echo.Context) error {
	return nil
}
