package courses

import (
	"context"
	"errors"

	materials "tourbackend/internal/courses/materials"
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/utils"
)

type Service struct {
	q                *db.Queries
	materialsService *materials.Service
}

func NewService(queries *db.Queries, materialsService *materials.Service) *Service {
	return &Service{
		queries,
		materialsService,
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
	Quizzes     []string             `json:"quizzes"`
}

func (s *Service) GetCourse(courseId string, host string, scheme string, ctx context.Context) (*GetCourseResponse, error) {

	course, err := s.q.GetCourse(ctx, courseId)
	if err != nil {
		if utils.IsNoRowsError(err) {
			return nil, CourseNotFound
		}
		return nil, FailedToFetchCourse
	}

	mats, err := s.materialsService.ListMaterials(courseId, host, scheme, ctx)

	courseDetail := GetCourseResponse{
		Uuid:        course.Uuid,
		Name:        course.Name,
		Description: course.Description,
		Materials:   mats,
		Quizzes:     []string{},
	}

	return &courseDetail, nil
}

func (s *Service) UpdateCourse(params db.UpdateCourseParams, ctx context.Context) (*db.Course, error) {

	updated, err := s.q.UpdateCourse(ctx, params)
	if err != nil {
		if utils.IsNoRowsError(err) {
			return nil, CourseNotFound
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
		return CourseNotFound
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
