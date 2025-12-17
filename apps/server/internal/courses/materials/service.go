package materials

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	database "tourbackend/internal/database/gen"

	"github.com/gabriel-vasile/mimetype"
)

var MAX_SIZE = int64(30 * 1024 * 1024)

var ALLOWED_FILES = map[string]bool{
	"application/pdf": true, // .pdf
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
	"text/plain": true, // .txt
	"image/png":  true, // .png
	"image/jpeg": true, // .jpeg and .jpg
	"image/gif":  true, // .gif
	"video/mp4":  true, // .mp4
	"audio/mpeg": true, // .mp3
}

var MIME_TO_EXT = map[string]string{
	"application/pdf": ".pdf",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	"text/plain": ".txt",
	"image/png":  ".png",
	"image/jpeg": ".jpeg and .jpg",
	"image/gif":  ".gif",
	"video/mp4":  ".mp4",
	"audio/mpeg": ".mp3",
}

type Material interface {
	GetType() string
}

type FileMaterial struct {
	Uuid        string `json:"uuid"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FileUrl     string `json:"fileUrl"`
	MimeType    string `json:"mimeType"`
	SizeBytes   int    `json:"sizeBytes"`
}

func (f FileMaterial) GetType() string {
	return f.Type
}

type UrlMaterial struct {
	Uuid        string `json:"uuid"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	FaviconUrl  string `json:"faviconUrl"`
}

func (f UrlMaterial) GetType() string {
	return f.Type
}

type Service struct {
	q          *database.Queries
	staticPath string
}

func NewService(queries *database.Queries, staticPath string) *Service {
	return &Service{queries, staticPath}
}

func (s *Service) deriveFaviconUrl(url string) string {
	if strings.HasPrefix(url, "http") {
		parts := strings.Split(url, "//")
		host := parts[1]
		base := strings.Split(host, "/")[0]
		faviconUrl := base + "/favicon.ico"
		fmt.Println(parts, host, base, faviconUrl)
		return faviconUrl
	}

	urlParts := strings.Split(url, "/")
	faviconUrl := urlParts[0] + "/favicon.ico"
	return faviconUrl
}

func (s *Service) ListMaterials(courseId string, host string, scheme string, ctx context.Context) ([]Material, error) {

	dbMaterials, err := s.q.ListAllMaterialsOfCourse(ctx, courseId)
	if err != nil {
		fmt.Println(err)
		return nil, FailedToFetchMaterials
	}

	if len(dbMaterials) == 0 {
		ok, err := s.q.CheckCourseExists(ctx, courseId)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("failed to check the course id")
		}

		if ok != 1 {
			return nil, CourseNotFound
		}
		return []Material{}, nil
	}

	formattedMaterials := make([]Material, 0, len(dbMaterials))
	for _, material := range dbMaterials {

		fmt.Println(material)

		expectedUrlStartForLocalFile := scheme + "://" + host + "/api/static/uploads"
		if strings.HasPrefix(material.Url, expectedUrlStartForLocalFile) {

			formattedMaterials = append(formattedMaterials, FileMaterial{
				Uuid:        material.Uuid,
				Type:        "file",
				Name:        material.Name,
				Description: material.Description,
				FileUrl:     material.Url,
				MimeType:    "mime",
				SizeBytes:   0,
			})
		} else {

			faviconUrl := s.deriveFaviconUrl(material.Url)

			formattedMaterials = append(formattedMaterials, UrlMaterial{
				Uuid:        material.Uuid,
				Type:        "url",
				Name:        material.Name,
				Description: material.Description,
				Url:         material.Url,
				FaviconUrl:  faviconUrl,
			})
		}
	}
	return formattedMaterials, nil
}

func (s *Service) checkFileSize(fileHeader *multipart.FileHeader) bool {
	if fileHeader.Size > MAX_SIZE {
		return false
	}
	return true
}

func (s *Service) checkFileType(src multipart.File) (bool, string, error) {

	m, err := mimetype.DetectReader(src)
	if err != nil {
		return false, "", err
	}

	mimeType := m.String()
	moreParts := strings.Split(mimeType, ";")
	if len(moreParts) > 1 {
		mimeType = moreParts[0]
	}
	fmt.Println(mimeType)

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return false, "", FailedToRewindCursorAfterFileTypeCheck
	}

	if !ALLOWED_FILES[mimeType] {
		return false, "", nil
	}

	return true, MIME_TO_EXT[mimeType], nil
}

func (s *Service) saveFileToStatic(
	params *CreateFileMaterialParams,
	src multipart.File,
	ext string,
) (string, error) {

	pathToUploads := s.staticPath + "/uploads/"
	pathToCourseFolder := pathToUploads + "/" + params.courseId
	pathToMaterialsFolder := pathToCourseFolder + "/materials"

	pathToFile := pathToMaterialsFolder + "/" + params.uuid + ext

	err := os.Mkdir(pathToUploads, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			fmt.Println("failed to create materials folder, wrong path", err)
			return "", err
		}
	}

	err = os.Mkdir(pathToCourseFolder, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			fmt.Println("failed to create course folder, wrong path", err)
			return "", err
		}
	}

	err = os.Mkdir(pathToMaterialsFolder, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			fmt.Println("failed to create materials folder, wrong path", err)
			return "", err
		}
	}

	dst, err := os.Create(pathToFile)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return params.scheme + "://" + params.host + "/api/static/uploads/" + params.courseId + "/materials/" + params.uuid + ext, nil

}

type CreateFileMaterialParams struct {
	fileHeader *multipart.FileHeader

	uuid        string
	name        string
	description string
	courseId    string

	scheme string
	host   string
}

func (s *Service) checkAndSaveFile(params *CreateFileMaterialParams) (string, error) {
	if !s.checkFileSize(params.fileHeader) {
		return "", TooBigMaterialFile
	}

	src, err := params.fileHeader.Open()
	if err != nil {
		fmt.Println(err)
		return "", FailedToOpenMaterialFile
	}
	defer src.Close()

	ok, ext, err := s.checkFileType(src)
	if err != nil {
		return "", FailedToRewindCursorAfterFileTypeCheck
	}
	if !ok {
		return "", ForbiddenFileType
	}

	url, err := s.saveFileToStatic(params, src, ext)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *Service) CreateFileMaterial(params *CreateFileMaterialParams, ctx context.Context) (*database.Material, error) {

	url, err := s.checkAndSaveFile(params)
	if err != nil {
		return nil, err
	}

	dbMaterial, err := s.q.CreateMaterial(ctx, database.CreateMaterialParams{
		Uuid:        params.uuid,
		Name:        params.name,
		Description: params.description,
		Url:         url,
		Courseuuid:  params.courseId,
	})
	if err != nil {
		return nil, err
	}

	return &dbMaterial, nil
}

type CreateUrlMaterialParams struct {
	uuid        string
	name        string
	description string
	courseId    string

	url string
}

func (s *Service) CreateUrlMaterial(params CreateUrlMaterialParams, ctx context.Context) (*database.Material, error) {

	dbMaterial, err := s.q.CreateMaterial(ctx, database.CreateMaterialParams{
		Uuid:        params.uuid,
		Name:        params.name,
		Description: params.description,
		Url:         params.url,
		Courseuuid:  params.courseId,
	})
	if err != nil {
		return nil, err
	}

	return &dbMaterial, nil
}

func (s *Service) DeleteMaterial(materialId string, ctx context.Context) error {

	res, err := s.q.DeleteMaterial(ctx, materialId)
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to delete material from db")
	}

	n, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to delete material from db")
	}

	if n == 0 {
		return CourseNotFound
	}

	return nil
}

func (s *Service) UpdateFileMaterial(params *CreateFileMaterialParams, ctx context.Context) (*database.Material, error) {

	url, err := s.checkAndSaveFile(params)
	if err != nil {
		return nil, err
	}

	material, err := s.q.UpdateMaterial(ctx, database.UpdateMaterialParams{
		Name:        params.name,
		Uuid:        params.uuid,
		Url:         url,
		Description: params.description,
	})
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (s *Service) UpdateUrlMaterial(params *CreateUrlMaterialParams, ctx context.Context) (*database.Material, error) {

	material, err := s.q.UpdateMaterial(ctx, database.UpdateMaterialParams{
		Name:        params.name,
		Uuid:        params.uuid,
		Url:         params.url,
		Description: params.description,
	})
	if err != nil {
		return nil, err
	}
	return &material, nil
}
