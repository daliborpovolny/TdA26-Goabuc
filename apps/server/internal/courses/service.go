package courses

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	materials "tourbackend/internal/courses/materials"
	"tourbackend/internal/courses/quizzes"
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/feeds"
	"tourbackend/internal/utils"
)

var ALLOWED_COURSE_STATES []string = []string{
	"preparation",
	"open",
	"closed",
}

var ALLOWED_MODULE_STATES []string = []string{
	"preparation",
	"open",
	"closed",
}

type Service struct {
	q                *db.Queries
	materialsService *materials.Service
	quizzesService   *quizzes.Service
	feedsService     *feeds.Service
}

func NewService(queries *db.Queries, materialsService *materials.Service, quizzesService *quizzes.Service, feedsService *feeds.Service) *Service {
	return &Service{
		queries,
		materialsService,
		quizzesService,
		feedsService,
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
	Uuid        string                   `json:"uuid"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Materials   []materials.Material     `json:"materials"`
	Quizzes     []quizzes.Quiz           `json:"quizzes"`
	Feed        []feeds.FeedPostResponse `json:"feed"`
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

	feed, err := s.feedsService.GetFeed(ctx, courseId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if feed == nil {
		feed = []feeds.FeedPostResponse{}
	}

	courseDetail := GetCourseResponse{
		Uuid:        course.Uuid,
		Name:        course.Name,
		Description: course.Description,
		Materials:   mats,
		Quizzes:     quizzes,
		Feed:        feed,
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

	s.feedsService.CreateAutomaticPost("Course updated", params.Uuid, ctx)

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

func (s *Service) ChangeCourseState(courseId string, state string, ctx context.Context) (db.Course, error) {

	if !slices.Contains(ALLOWED_COURSE_STATES, state) {
		return db.Course{}, ErrBadCourseState
	}

	now := time.Now().Unix()

	course, err := s.q.ChangeCourseState(ctx, db.ChangeCourseStateParams{
		State:     state,
		Uuid:      courseId,
		UpdatedAt: now,
	})
	if err != nil {
		return db.Course{}, err
	}

	return course, err
}

// modules

func (s *Service) CreateModule(courseId string, moduleId string, name string, ctx context.Context) (db.Module, error) {

	now := time.Now().Unix()

	module, err := s.q.CreateModule(ctx, db.CreateModuleParams{
		Uuid:       moduleId,
		CourseUuid: courseId,
		Name:       name,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
	if err != nil {
		return db.Module{}, err
	}

	return module, nil
}

func (s *Service) ChangeModuleState(courseId string, moduleId string, state string, ctx context.Context) (db.Module, error) {

	if !slices.Contains(ALLOWED_MODULE_STATES, state) {
		return db.Module{}, ErrBadModuleState
	}

	now := time.Now().Unix()

	module, err := s.q.ChangeModuleState(ctx, db.ChangeModuleStateParams{
		State:     state,
		Uuid:      moduleId,
		UpdatedAt: now,
	})
	if err != nil {
		return db.Module{}, err
	}

	return module, err
}

func (s *Service) GetModule(courseId string, moduleId string, ctx context.Context) (db.Module, error) {

	module, err := s.q.GetModule(ctx, db.GetModuleParams{
		Uuid:       moduleId,
		CourseUuid: courseId,
	})
	if err != nil {
		return db.Module{}, err
	}

	return module, nil
}

func (s *Service) UpdateModule(courseId string, moduleId string, name string, ctx context.Context) (db.Module, error) {

	newModule, err := s.q.UpdateModule(ctx, db.UpdateModuleParams{
		Uuid:       moduleId,
		CourseUuid: courseId,
		Name:       name,
	})
	if err != nil {
		return db.Module{}, err
	}

	return newModule, nil
}
