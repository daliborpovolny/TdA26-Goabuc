package courses

import "errors"

var (
	ErrCourseNotFound      = errors.New("No course with such id exists")
	ErrFailedToFetchCourse = errors.New("Failed to fetch course from db")
	ErrBadCourseState      = errors.New("Invalid new course state")
	ErrBadModuleState      = errors.New("Invalid new module state")
)
