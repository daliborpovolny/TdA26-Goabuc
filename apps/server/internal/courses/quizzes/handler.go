package quizzes

import (
	"net/http"
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

func (h *Handler) ListQuizzes(c echo.Context) error {
	r := h.NewReqCtx(c)

	quizzes, err := h.service.ListQuizes(r.Ctx)
	if err != nil {
		return r.ServerError(err)
	}
	return c.JSON(http.StatusOK, quizzes)
}

func (h *Handler) CreateQuizz(c echo.Context) error {
	r := h.NewReqCtx(c)

	courseId := r.Echo.Param("courseId")

	var quiz Quiz
	if err := c.Bind(&quiz); err != nil {
		return r.Error(http.StatusBadRequest, "bad request")
	}

	if quiz.Uuid == "" {
		quiz.Uuid = uuid.NewString()
	}

	for i := range quiz.Questions {
		if quiz.Questions[i].Uuid == "" {
			quiz.Questions[i].Uuid = uuid.NewString()
		}
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
	r := h.NewReqCtx(c)

	quizId := r.Echo.Param("quizId")

	quiz, err := h.service.q.GetQuiz(r.Ctx, quizId)
	if err != nil {
		if err == ErrQuizNotFound {
			return r.Error(http.StatusBadRequest, "unknown quiz id")
		}
		return r.ServerError(err)
	}
	return c.JSON(http.StatusOK, quiz)
}

func (h *Handler) UpdateQuizz(c echo.Context) error {
	// r := h.NewReqCtx(c)

	// courseId := r.Echo.Param("courseId")

	// quiz, err := h.service.UpdateQuiz(, r.Ctx)
	// if err != nil {
	// 	if err == ErrQuizNotFound {
	// 		return r.Error(http.StatusOK, "unknown quiz id")
	// 	}
	// 	if err == ErrBadQuestionType {
	// 		return r.Error(http.StatusBadRequest, "invalid question type")
	// 	}
	// 	return r.ServerError(err)
	// }

	// return c.JSON(http.StatusCreated, quiz)
	return nil

}

func (h *Handler) DeleteQuizz(c echo.Context) error {
	r := h.NewReqCtx(c)

	quizId := r.Echo.Param("quizId")
	err := h.service.DeleteQuiz(quizId, r.Ctx)
	if err != nil {
		if err == ErrQuizNotFound {
			return r.Error(http.StatusNotFound, "quiz not found")
		}
		return r.ServerError(err)
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) SubmitQuizAnswers(c echo.Context) error {
	return nil
}
