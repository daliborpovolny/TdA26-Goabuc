package quizzes

import "errors"

var (
	ErrQuizNotFound       = errors.New("Quiz not found")
	ErrBadQuestionType    = errors.New("Question type can be either singleChoice or multipleChoice")
	ErrBadNumberOfAnswers = errors.New("Number of answers must match the number of questions")
)
