package materials

import "errors"

var (
	FailedToFetchMaterials                 = errors.New("failed to fetch materials from db")
	FailedToOpenMaterialFile               = errors.New("failed to open the material file")
	FailedToRewindCursorAfterFileTypeCheck = errors.New("failed to rewind cursor after file type check")
	TooBigMaterialFile                     = errors.New("too big material file max is 30MB")
	ForbiddenFileType                      = errors.New("forbidden file type")
	CourseNotFound                         = errors.New("unknown course id")
)
