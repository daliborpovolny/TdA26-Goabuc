package materials

import "errors"

var (
	ErrFileTooBig        = errors.New("too big material file max is 30MB")
	ErrFileTypeForbidden = errors.New("forbidden file type")
	ErrCourseNotFound    = errors.New("unknown course id")
)
