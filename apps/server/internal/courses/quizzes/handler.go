package quizzes

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

	courseId := c.Param("courseId")

	quizzes, err := h.service.ListQuizes(courseId, r.Ctx)
	if err != nil {
		return r.ServerError(err)
	}
	return c.JSON(http.StatusOK, quizzes)
}

func (h *Handler) CreateQuiz(c echo.Context) error {
	r := h.NewReqCtx(c)

	courseId := r.Echo.Param("courseId")

	var quiz Quiz
	if err := c.Bind(&quiz); err != nil {
		return r.Error(http.StatusBadRequest, "bad request")
	}

	quiz.Uuid = uuid.NewString()

	if len(quiz.Questions) < 1 {
		return r.Error(http.StatusBadRequest, "quiz must have at least one question")
	}

	for i := range quiz.Questions {
		if quiz.Questions[i].Uuid == "" {
			quiz.Questions[i].Uuid = uuid.NewString()
		}
		// quiz.Questions[i].Uuid = uuid.NewString()
	}

	dbQuiz, err := h.service.CreateQuiz(quiz, courseId, r.Ctx)
	if err != nil {
		var eqbf *ErrQuestionBadFormat

		if errors.As(err, &eqbf) {
			fmt.Println("bad format sucker!")
			return r.Error(http.StatusBadRequest, eqbf.Error())
		}

		return r.ServerError(err)
	}
	return r.Echo.JSON(http.StatusCreated, dbQuiz)
}

func (h *Handler) GetQuiz(c echo.Context) error {
	r := h.NewReqCtx(c)

	quizId := r.Echo.Param("quizId")

	quiz, err := h.service.GetQuiz(quizId, r.Ctx)
	if err != nil {
		if err == ErrQuizNotFound {
			return r.Error(http.StatusNotFound, "unknown quiz id")
		}
		return r.ServerError(err)
	}
	return c.JSON(http.StatusOK, quiz)
}

func (h *Handler) UpdateQuiz(c echo.Context) error {
	r := h.NewReqCtx(c)

	quizId := r.Echo.Param("quizId")

	var quiz Quiz
	if err := c.Bind(&quiz); err != nil {
		return r.Error(http.StatusBadRequest, "bad request")
	}

	if len(quiz.Questions) < 1 {
		return r.Error(http.StatusBadRequest, "quiz must have at least one question")
	}

	for i := range quiz.Questions {
		if quiz.Questions[i].Uuid == "" {
			quiz.Questions[i].Uuid = uuid.NewString()
		}
	}

	quiz.Uuid = quizId

	updatedQuiz, err := h.service.UpdateQuiz(&quiz, r.Ctx)
	if err != nil {
		if err == ErrQuizNotFound {
			return r.Error(http.StatusOK, "unknown quiz id")
		}
		if err == ErrBadQuestionType {
			return r.Error(http.StatusBadRequest, "invalid question type")
		}

		var eqbf *ErrQuestionBadFormat

		if errors.As(err, &eqbf) {
			fmt.Println("bad format sucker!")
			return r.Error(http.StatusBadRequest, eqbf.Error())
		}

		return r.ServerError(err)
	}

	return c.JSON(http.StatusCreated, updatedQuiz)
}

func (h *Handler) DeleteQuiz(c echo.Context) error {
	r := h.NewReqCtx(c)

	quizId := r.Echo.Param("quizId")
	err := h.service.DeleteQuiz(quizId, r.Ctx)
	if err != nil {
		if err == ErrQuizNotFound {
			return r.Error(http.StatusNotFound, "quiz not found")
		}
		return r.ServerError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

type Answer struct {
	Uuid    string `json:"uuid"`
	Comment string `json:"comment"`

	// singleChoice
	SelectedIndex *int `json:"selectedIndex,omitempty"`

	// multipleChoice
	SelectedIndices []int `json:"selectedIndices,omitempty"`
}

type SubmitQuizAnswersRequest struct {
	Comment string   `json:"comment"`
	Answers []Answer `json:"answers"`
	UserID  *int     `json:"userId"`
}

type SubmittedAnswersOutcome struct {
	QuizUuid           string `json:"quizUuid"`
	Score              int    `json:"score"`
	MaxScore           int    `json:"maxScore"`
	CorrectPerQuestion []bool `json:"correctPerQuestion"`
	SubmittedAt        string `json:"submittedAt"`
}

func (h *Handler) SubmitQuizAnswers(c echo.Context) error {
	r := h.NewReqCtx(c)

	quizId := r.Echo.Param("quizId")

	var answers SubmitQuizAnswersRequest
	if err := c.Bind(&answers); err != nil {
		return r.Error(http.StatusBadRequest, "bad request")
	}

	outcome, err := h.service.SubmitQuizAnswers(quizId, answers, r.Ctx)
	if err != nil {
		if err == ErrBadNumberOfAnswers {
			return r.Error(http.StatusBadRequest, err.Error())
		}

		var ebr *ErrBadRequest
		if errors.As(err, &ebr) {
			return r.Error(http.StatusBadRequest, ebr.Error())
		}

		fmt.Println("ERRORROOOOO: \n\n\n\n", err)
		return r.ServerError(err)
	}

	return c.JSON(http.StatusOK, outcome)
}

func (h *Handler) GetAnswersOfQuiz(c echo.Context) error {
	r := h.NewReqCtx(c)

	fmt.Println("in the handler")

	quizId := r.Echo.Param("quizId")
	if quizId == "" {
		return r.Error(http.StatusBadRequest, "Quizid must be set")
	}

	answers, err := h.service.GetAnswersOfQuiz(quizId, r.Ctx)
	if err != nil {
		return r.ServerError(err)
	}

	fmt.Println(answers)

	return c.JSON(http.StatusOK, answers)
}
