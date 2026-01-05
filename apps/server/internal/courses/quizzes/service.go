package quizzes

import (
	"context"
	"fmt"
	"time"
	db "tourbackend/internal/database/gen"

	"github.com/google/uuid"
)

type Service struct {
	q          *db.Queries
	staticPath string
}

func NewService(queries *db.Queries, staticPath string) *Service {
	fmt.Println("quizes serevice!")
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
		dbQuestion, err := s.q.CreateQuestion(ctx, db.CreateQuestionParams{
			Uuid:     uuid.NewString(),
			QuizUuid: quiz.Uuid,
			Type:     question.QueType,
		})
	}

	return s.dbQuizToQuiz(dbQuiz)

}

func (s *Service) dbQuizToQuiz(dbQuiz db.Quizz, questions []db.Question) *Quiz {

}

func (s *Service) UpdateQuiz() (*Quiz, error) {
	return nil, nil
}

func (s *Service) GetQuiz() (*Quiz, error) {
	return nil, nil
}

func (s *Service) ListQuizes() ([]Quiz, error) {
	return nil, nil
}

func (s *Service) DeleteQuiz() error {
	return nil
}

func (s *Service) DeleteQuestion() error {
	return nil
}
