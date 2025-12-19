package courses

import "errors"

var (
	ErrCourseNotFound      = errors.New("No course with such id exists")
	ErrFailedToFetchCourse = errors.New("Failed to fetch course from db")
)
