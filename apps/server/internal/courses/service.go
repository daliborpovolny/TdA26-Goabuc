package courses

import (
	"cmp"
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

type Module struct {
	Uuid       string `json:"uuid"`
	CourseUuid string `json:"courseUuid"`

	Name        string `json:"name"`
	Description string `json:"description"`
	State       string `json:"state"` // one of ALLOWED_MODULE_STATES

	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	// Items        []Item `json:"items"`
	// NewItemOrder int    `json:"newItemOrder"`
}

type FullModule struct {
	Module
	Items        []Item `json:"items"`
	NewItemOrder int    `json:"newItemOrder"`
}

func (s *Service) dbModuleToModule(dbM db.Module) Module {
	return Module{
		Uuid:       dbM.Uuid,
		CourseUuid: dbM.CourseUuid,

		Name:        dbM.Name,
		Description: dbM.Description,
		State:       dbM.State,

		CreatedAt: utils.UnixToIso(dbM.CreatedAt),
		UpdatedAt: utils.UnixToIso(dbM.UpdatedAt),
	}
}

func (s *Service) moduleToFullModule(m Module, items []Item, newItemOrder int) FullModule {
	return FullModule{
		Module:       m,
		Items:        items,
		NewItemOrder: newItemOrder,
	}
}

type Item interface {
	GetModuleOrder() int
	GetModuleId() string
}

func (s *Service) CreateCourse(params db.CreateCourseParams, ctx context.Context) (*db.Course, error) {
	course, err := s.q.CreateCourse(ctx, params)
	if err != nil {
		return nil, err
	}

	return &course, nil
}

type GetCourseResponse struct {
	Uuid string `json:"uuid"`

	Name        string `json:"name"`
	Description string `json:"description"`
	State       string `json:"state"`
	Archived    bool   `json:"archived"`

	Materials []materials.Material `json:"materials"`
	Quizzes   []quizzes.Quiz       `json:"quizzes"`

	Feed []feeds.FeedPostResponse `json:"feed"`

	Modules []FullModule `json:"modules"`
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

	modules, err := s.ListAllModules(courseId, ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fullModules := make([]FullModule, 0, len(modules))

	for _, module := range modules {
		items := make([]Item, 0, 10)
		maxItemOrder := 0

		for _, quiz := range quizzes {
			if quiz.ModuleId == module.Uuid {
				items = append(items, quiz)
				maxItemOrder = max(quiz.GetModuleOrder(), maxItemOrder)
			}
		}

		for _, mat := range mats {
			if mat.GetModuleId() == module.Uuid {
				items = append(items, mat)
				maxItemOrder = max(mat.GetModuleOrder(), maxItemOrder)

			}
		}

		slices.SortFunc(items, func(a, b Item) int {
			return cmp.Compare(a.GetModuleOrder(), b.GetModuleOrder())
		})

		fullModules = append(fullModules, s.moduleToFullModule(module, items, maxItemOrder+1))

	}

	if feed == nil {
		feed = []feeds.FeedPostResponse{}
	}

	courseDetail := GetCourseResponse{
		Uuid: course.Uuid,

		Name:        course.Name,
		Description: course.Description,
		State:       course.State,
		Archived:    course.Archived == 1,

		Materials: mats,
		Quizzes:   quizzes,

		Feed: feed,

		Modules: fullModules,
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

//* Modules

func (s *Service) CreateModule(courseId string, moduleId string, name string, description string, ctx context.Context) (Module, error) {

	now := time.Now().Unix()

	dbModule, err := s.q.CreateModule(ctx, db.CreateModuleParams{
		Uuid:        moduleId,
		CourseUuid:  courseId,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return Module{}, err
	}

	return s.dbModuleToModule(dbModule), nil
}

func (s *Service) ChangeModuleState(courseId string, moduleId string, state string, ctx context.Context) (Module, error) {

	if !slices.Contains(ALLOWED_MODULE_STATES, state) {
		return Module{}, ErrBadModuleState
	}

	now := time.Now().Unix()

	dbModule, err := s.q.ChangeModuleState(ctx, db.ChangeModuleStateParams{
		State:     state,
		Uuid:      moduleId,
		UpdatedAt: now,
	})
	if err != nil {
		return Module{}, err
	}

	return s.dbModuleToModule(dbModule), err
}

func (s *Service) GetModule(courseId string, moduleId string, ctx context.Context) (Module, error) {

	module, err := s.q.GetModule(ctx, db.GetModuleParams{
		Uuid:       moduleId,
		CourseUuid: courseId,
	})
	if err != nil {
		return Module{}, err
	}

	return s.dbModuleToModule(module), nil
}

func (s *Service) ListAllModules(courseId string, ctx context.Context) ([]Module, error) {

	dbModules, err := s.q.ListAllModules(ctx, courseId)
	if err != nil {
		return nil, err
	}

	modules := make([]Module, 0, len(dbModules))
	for _, dbM := range dbModules {
		modules = append(modules, s.dbModuleToModule(dbM))
	}

	return modules, nil
}

func (s *Service) UpdateModule(courseId string, moduleId string, name string, description string, ctx context.Context) (Module, error) {

	newModule, err := s.q.UpdateModule(ctx, db.UpdateModuleParams{
		Uuid:        moduleId,
		CourseUuid:  courseId,
		Name:        name,
		Description: description,
	})
	if err != nil {
		return Module{}, err
	}

	return s.dbModuleToModule(newModule), nil
}

func (s *Service) DeleteModule(courseId string, moduleId string, ctx context.Context) error {

	err := s.q.DeleteModule(ctx, db.DeleteModuleParams{
		CourseUuid: courseId,
		Uuid:       moduleId,
	})
	if err != nil {
		return err
	}
	return nil
}
