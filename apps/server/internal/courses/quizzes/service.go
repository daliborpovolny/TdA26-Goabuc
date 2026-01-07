package quizzes

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/utils"

	"github.com/google/uuid"
)

type Service struct {
	q          *db.Queries
	staticPath string
}

func NewService(queries *db.Queries, staticPath string) *Service {
	return &Service{queries, staticPath}
}

type Quiz struct {
	Uuid          string     `json:"uuid"`
	Title         string     `json:"title"`
	AttemptsCount int        `json:"attemptsCount"`
	Questions     []Question `json:"questions"`
}

type Question struct {
	Uuid    string `json:"uuid"`
	QueType string `json:"type"` // "singleChoice" | "multipleChoice"

	Question string   `json:"question"`
	Options  []string `json:"options"`

	// singleChoice
	CorrectIndex *int `json:"correctIndex,omitempty"`

	// multipleChoice
	CorrectIndices []int `json:"correctIndices,omitempty"`
}

func (s *Service) validateQuestions(questions []Question) bool {
	for _, q := range questions {

		if uuid.Validate(q.Uuid) != nil || q.Uuid == "" {
			return false
		}

		switch q.QueType {
		case "singleChoice":
			if q.CorrectIndex == nil {
				return false
			}
		case "multipleChoice":
			if len(q.CorrectIndices) == 0 {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (s *Service) CreateQuiz(quiz Quiz, courseId string, ctx context.Context) (*Quiz, error) {

	now := time.Now().Unix()

	dbQuiz, err := s.q.CreateQuizz(ctx, db.CreateQuizzParams{
		Uuid:          quiz.Uuid,
		CourseUuid:    courseId,
		Title:         quiz.Title,
		AttemptsCount: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		return nil, err
	}

	dbQuestions := make([]db.Question, 0, len(quiz.Questions))
	for _, question := range quiz.Questions {

		var stringIndices []string
		switch question.QueType {
		case "singleChoice":
			stringIndices = []string{strconv.Itoa(*question.CorrectIndex)}
		case "multipleChoice":
			stringIndices = make([]string, 0, len(question.CorrectIndices))
			for _, index := range question.CorrectIndices {
				stringIndices = append(stringIndices, strconv.Itoa(index))
			}
		}

		dbQuestion, err := s.q.CreateQuestion(ctx, db.CreateQuestionParams{
			Uuid:           uuid.NewString(),
			QuizUuid:       quiz.Uuid,
			Type:           question.QueType,
			QuestionText:   question.Question,
			Options:        strings.Join(question.Options, "|"),
			CorrectIndices: strings.Join(stringIndices, "|"),
		})
		if err != nil {
			return nil, err
		}

		dbQuestions = append(dbQuestions, dbQuestion)
	}

	return s.dbQuizToQuiz(dbQuiz, dbQuestions)

}

func (s *Service) dbQuizToQuiz(dbQuiz db.Quizz, questions []db.Question) (*Quiz, error) {
	quiz := &Quiz{
		Uuid:          dbQuiz.Uuid,
		Title:         dbQuiz.Title,
		AttemptsCount: int(dbQuiz.AttemptsCount),
	}
	quiz.Questions = make([]Question, 0, len(questions))

	for _, dbQue := range questions {
		question, err := s.dbQuestionToQuestion(dbQue)
		if err != nil {
			return nil, err
		}

		quiz.Questions = append(quiz.Questions, question)
	}

	return quiz, nil
}

func (s *Service) dbQuestionToQuestion(dbQue db.Question) (Question, error) {

	options := strings.Split(dbQue.Options, "|")

	correctStringIndices := strings.Split(dbQue.CorrectIndices, "|")
	correctIndices := make([]int, 0, len(correctStringIndices))
	for _, stringIndex := range correctStringIndices {
		index, err := strconv.Atoi(stringIndex)
		if err != nil {
			return Question{}, err
		}
		correctIndices = append(correctIndices, index)
	}

	question := Question{
		Uuid:    dbQue.Uuid,
		QueType: dbQue.Type,

		Question: dbQue.QuestionText,
		Options:  options,
	}

	switch dbQue.Type {
	case "singleChoice":
		question.CorrectIndex = &correctIndices[0]
	case "multipleChoice":
		question.CorrectIndices = correctIndices
	default:
		return Question{}, errors.New("unknown question type")
	}

	return question, nil

}

func (s *Service) UpdateQuiz(quiz *Quiz, ctx context.Context) (*Quiz, error) {

	dbQuiz, err := s.q.UpdateQuizz(ctx, db.UpdateQuizzParams{
		Title:         utils.ToSqlNullString(&quiz.Title),
		AttemptsCount: sql.NullInt64{Int64: 0, Valid: false},
		UpdatedAt:     sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
		Uuid:          quiz.Uuid,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.q.DeleteQuestionsOfQuiz(ctx, quiz.Uuid)
	if err != nil {
		return nil, err
	}

	dbQuestions := make([]db.Question, 0, len(quiz.Questions))
	for _, question := range quiz.Questions {

		var stringIndices []string
		switch question.QueType {
		case "singleChoice":
			stringIndices = []string{strconv.Itoa(*question.CorrectIndex)}
		case "multipleChoice":
			stringIndices = make([]string, 0, len(question.CorrectIndices))
			for _, index := range question.CorrectIndices {
				stringIndices = append(stringIndices, strconv.Itoa(index))
			}
		}

		dbQuestion, err := s.q.CreateQuestion(ctx, db.CreateQuestionParams{
			Uuid:           uuid.NewString(),
			QuizUuid:       quiz.Uuid,
			Type:           question.QueType,
			QuestionText:   question.Question,
			Options:        strings.Join(question.Options, "|"),
			CorrectIndices: strings.Join(stringIndices, "|"),
		})
		if err != nil {
			return nil, err
		}

		dbQuestions = append(dbQuestions, dbQuestion)
	}

	return s.dbQuizToQuiz(dbQuiz, dbQuestions)
}

func (s *Service) convertGetQuizRowsToQuiz(rows []db.GetQuizRow) (*Quiz, error) {
	if len(rows) < 1 {
		return nil, ErrQuizNotFound
	}

	r := rows[0]

	quiz := &Quiz{
		Uuid:          r.QuizUuid,
		Title:         r.QuizTitle,
		AttemptsCount: int(r.QuizAttemptsCount),
		Questions:     make([]Question, 0, len(rows)),
	}

	for _, qr := range rows {

		options := strings.Split(qr.QuestionOptions.String, "|")

		correctStringIndices := strings.Split(qr.QuestionCorrectIndices.String, "|")
		correctIndices := make([]int, 0, len(correctStringIndices))
		for _, stringIndex := range correctStringIndices {
			index, err := strconv.Atoi(stringIndex)
			if err != nil {
				return nil, err
			}
			correctIndices = append(correctIndices, index)
		}
		qs := Question{
			Uuid:     qr.QuestionUuid.String,
			QueType:  qr.QuestionType.String,
			Question: qr.QuestionText.String,
			Options:  options,
		}

		switch qs.QueType {
		case "singleChoice":
			qs.CorrectIndex = &correctIndices[0]
		case "multipleChoice":
			qs.CorrectIndices = correctIndices
		default:
			return nil, ErrBadQuestionType
		}

		quiz.Questions = append(quiz.Questions, qs)
	}

	return quiz, nil
}

func (s *Service) GetQuiz(quizId string, ctx context.Context) (*Quiz, error) {

	rows, err := s.q.GetQuiz(ctx, quizId)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrQuizNotFound
	}

	quiz, err := s.convertGetQuizRowsToQuiz(rows)
	if err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *Service) convertListQuizRowsToQuizzes(rows []db.ListQuizzesRow) ([]Quiz, error) {
	if len(rows) < 1 {
		return []Quiz{}, nil
	}

	quizzes := make([]Quiz, 0, 5)
	var currentQuizUuid string
	currentQuizIndex := -1

	for _, qr := range rows {
		if currentQuizUuid != qr.QuizUuid {
			currentQuizIndex += 1
			quizzes = append(quizzes, Quiz{
				Uuid:          qr.QuizUuid,
				Title:         qr.QuizTitle,
				AttemptsCount: int(qr.QuizAttemptsCount),
				Questions:     make([]Question, 0, len(rows)),
			})
		}

		options := strings.Split(qr.QuestionOptions.String, "|")

		correctStringIndices := strings.Split(qr.QuestionCorrectIndices.String, "|")
		correctIndices := make([]int, 0, len(correctStringIndices))
		for _, stringIndex := range correctStringIndices {
			index, err := strconv.Atoi(stringIndex)
			if err != nil {
				return nil, err
			}
			correctIndices = append(correctIndices, index)
		}
		qs := Question{
			Uuid:     qr.QuestionUuid.String,
			QueType:  qr.QuestionType.String,
			Question: qr.QuestionText.String,
			Options:  options,
		}

		switch qs.QueType {
		case "singleChoice":
			qs.CorrectIndex = &correctIndices[0]
		case "multipleChoice":
			qs.CorrectIndices = correctIndices
		default:
			return nil, ErrBadQuestionType
		}

		quizzes[currentQuizIndex].Questions = append(quizzes[currentQuizIndex].Questions, qs)
	}
	return quizzes, nil
}

func (s *Service) ListQuizes(ctx context.Context) ([]Quiz, error) {

	rows, err := s.q.ListQuizzes(ctx)
	if err != nil {
		return nil, err
	}

	quizzes, err := s.convertListQuizRowsToQuizzes(rows)
	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (s *Service) DeleteQuiz(quizId string, ctx context.Context) error {

	res, err := s.q.DeleteQuizz(ctx, quizId)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return ErrQuizNotFound
	}

	return nil
}

func (s *Service) SubmitQuizAnswers() error {
	return nil
}
