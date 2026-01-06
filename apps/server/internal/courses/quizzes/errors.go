package quizzes

import "errors"

var (
	ErrQuizNotFound    = errors.New("Quiz not found")
	ErrBadQuestionType = errors.New("Question type can be either singleChoice or multipleChoice")
)
