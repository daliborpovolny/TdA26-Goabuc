package courses

import (
	"context"
	"errors"
	"fmt"

	materials "tourbackend/internal/courses/materials"
	"tourbackend/internal/courses/quizzes"
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/utils"
)

type Service struct {
	q                *db.Queries
	materialsService *materials.Service
	quizzesService   *quizzes.Service
}

func NewService(queries *db.Queries, materialsService *materials.Service, quizzesService *quizzes.Service) *Service {
	return &Service{
		queries,
		materialsService,
		quizzesService,
	}
}

func (s *Service) CreateCourse(params db.CreateCourseParams, ctx context.Context) (*db.Course, error) {
	course, err := s.q.CreateCourse(ctx, params)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

type GetCourseResponse struct {
	Uuid        string               `json:"uuid"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Materials   []materials.Material `json:"materials"`
	Quizzes     []quizzes.Quiz       `json:"quizzes"`
}

func (s *Service) GetCourse(courseId string, host string, scheme string, ctx context.Context) (*GetCourseResponse, error) {

	course, err := s.q.GetCourse(ctx, courseId)
	if err != nil {
		if utils.IsNoRowsError(err) {
			return nil, ErrCourseNotFound
		}
		return nil, ErrFailedToFetchCourse
	}

	mats, err := s.materialsService.ListMaterials(courseId, host, scheme, ctx)
	if err != nil {
		return nil, err
	}

	quizzes, err := s.quizzesService.ListQuizes(courseId, ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	courseDetail := GetCourseResponse{
		Uuid:        course.Uuid,
		Name:        course.Name,
		Description: course.Description,
		Materials:   mats,
		Quizzes:     quizzes,
	}

	return &courseDetail, nil
}

func (s *Service) UpdateCourse(params db.UpdateCourseParams, ctx context.Context) (*db.Course, error) {

	updated, err := s.q.UpdateCourse(ctx, params)
	if err != nil {
		if utils.IsNoRowsError(err) {
			return nil, ErrCourseNotFound
		}
		return nil, errors.New("Failed to update the course")
	}

	return &updated, nil
}

func (s *Service) DeleteCourse(courseId string, ctx context.Context) error {

	res, err := s.q.DeleteCourse(ctx, courseId)
	if err != nil {
		return errors.New("Failed to delete the course")
	}

	n, err := res.RowsAffected()
	if err != nil {
		return errors.New("Failed to delete the course")
	}

	if n == 0 {
		return ErrCourseNotFound
	}

	return nil
}

func (s *Service) ListAllCourses(ctx context.Context) ([]db.Course, error) {
	courses, err := s.q.ListAllCourses(ctx)
	if err != nil {
		return nil, errors.New("Failed to fetch courses from db")
	}
	return courses, nil
}
