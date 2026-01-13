package quizzes

import (
	"errors"
	"fmt"
)

var (
	ErrQuizNotFound       = errors.New("Quiz not found")
	ErrBadQuestionType    = errors.New("Question type can be either singleChoice or multipleChoice")
	ErrBadNumberOfAnswers = errors.New("Number of answers must match the number of questions")
)

type ErrQuestionBadFormat struct {
	questionNumber int
	message        string
}

func (e *ErrQuestionBadFormat) Error() string {
	return fmt.Sprintf(
		"question %d is wrong: %q",
		e.questionNumber,
		e.message,
	)
}

type ErrBadRequest struct {
	message string
}

func (e *ErrBadRequest) Error() string {
	return e.message
}
