package courses

import "errors"

var (
	CourseNotFound      = errors.New("No course with such id exists")
	FailedToFetchCourse = errors.New("Failed to fetch course from db")
)
