package materials

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	database "tourbackend/internal/database/gen"
	"tourbackend/internal/utils"

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

		expectedUrlStartForLocalFile := scheme + "://" + host + "/api/static/uploads"
		if strings.HasPrefix(material.Url, expectedUrlStartForLocalFile) {

			fileInfo, err := s.q.GetFileMaterialMetadata(ctx, material.Uuid)
			if err != nil {
				return nil, err
			}

			formattedMaterials = append(formattedMaterials, FileMaterial{
				Uuid:        material.Uuid,
				Type:        "file",
				Name:        material.Name,
				Description: material.Description,
				FileUrl:     material.Url,
				MimeType:    fileInfo.Mime,
				SizeBytes:   int(fileInfo.Size),
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

func (s *Service) checkFileSize(fileHeader *multipart.FileHeader) (bool, int64) {
	if fileHeader.Size > MAX_SIZE {
		return false, 0
	}
	return true, fileHeader.Size
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

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return false, "", FailedToRewindCursorAfterFileTypeCheck
	}

	if !ALLOWED_FILES[mimeType] {
		return false, "", nil
	}

	return true, mimeType, nil
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

	// delete former materials if any - this could be done only on updates but whatever!
	prefix := params.uuid
	entries, err := os.ReadDir(pathToMaterialsFolder)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		fmt.Println(name)

		if strings.HasPrefix(name, prefix) {

			err := os.Remove(filepath.Join(pathToMaterialsFolder, name))
			if err != nil {
				return "", err
			}
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

type SavedFileInfo struct {
	url  string
	size int64
	mime string
}

func (s *Service) checkAndSaveFile(params *CreateFileMaterialParams) (SavedFileInfo, error) {
	ok, size := s.checkFileSize(params.fileHeader)
	if !ok {
		return SavedFileInfo{}, TooBigMaterialFile
	}

	src, err := params.fileHeader.Open()
	if err != nil {
		fmt.Println(err)
		return SavedFileInfo{}, FailedToOpenMaterialFile
	}
	defer src.Close()

	ok, mimeType, err := s.checkFileType(src)
	if err != nil {
		return SavedFileInfo{}, FailedToRewindCursorAfterFileTypeCheck
	}
	if !ok {
		return SavedFileInfo{}, ForbiddenFileType
	}

	url, err := s.saveFileToStatic(params, src, MIME_TO_EXT[mimeType])
	if err != nil {
		return SavedFileInfo{}, err
	}

	return SavedFileInfo{url, size, mimeType}, nil
}

func (s *Service) CreateFileMaterial(params *CreateFileMaterialParams, ctx context.Context) (*database.Material, SavedFileInfo, error) {

	saveFileInfo, err := s.checkAndSaveFile(params)
	if err != nil {
		return nil, SavedFileInfo{}, err
	}

	dbMaterial, err := s.q.CreateMaterial(ctx, database.CreateMaterialParams{
		Uuid:        params.uuid,
		Name:        params.name,
		Description: params.description,
		Url:         saveFileInfo.url,
		Courseuuid:  params.courseId,
	})
	if err != nil {
		return nil, SavedFileInfo{}, err
	}

	_, err = s.q.CreateFileMaterialMetadata(ctx, database.CreateFileMaterialMetadataParams{
		MaterialUuid: dbMaterial.Uuid,
		Size:         saveFileInfo.size,
		Mime:         saveFileInfo.mime,
	})

	return &dbMaterial, saveFileInfo, nil
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

type UpdateFileMaterialParms struct {
	uuid     string
	courseId string

	fileHeader *multipart.FileHeader
	url        *string

	name        *string
	description *string

	scheme string
	host   string
}

func (s *Service) UpdateFileMaterial(params *UpdateFileMaterialParms, ctx context.Context) (*database.Material, SavedFileInfo, error) {

	if params.fileHeader != nil {

		saveFileInfo, err := s.checkAndSaveFile(&CreateFileMaterialParams{
			fileHeader:  params.fileHeader,
			uuid:        params.uuid,
			courseId:    params.courseId,
			name:        "",
			description: "",
			scheme:      params.scheme,
			host:        params.host,
		})
		if err != nil {
			return nil, SavedFileInfo{}, err
		}

		_, err = s.q.UpdateFileMaterialMetadata(ctx, database.UpdateFileMaterialMetadataParams{
			MaterialUuid: params.uuid,
			Size:         saveFileInfo.size,
			Mime:         saveFileInfo.mime,
		})
		if err != nil {
			return nil, SavedFileInfo{}, err
		}

		params.url = &saveFileInfo.url
	}

	material, err := s.q.UpdateMaterialPartial(ctx, database.UpdateMaterialPartialParams{
		Name:        utils.ToSqlNullString(params.name),
		Uuid:        params.uuid,
		Url:         utils.ToSqlNullString(params.url),
		Description: utils.ToSqlNullString(params.description),
	})
	if err != nil {
		return nil, SavedFileInfo{}, err
	}

	dbInfo, err := s.q.GetFileMaterialMetadata(ctx, params.uuid)

	return &material, SavedFileInfo{material.Url, dbInfo.Size, dbInfo.Mime}, nil
}

func (s *Service) UpdateUrlMaterial(name *string, description *string, url *string, uuid string, ctx context.Context) (*database.Material, error) {

	material, err := s.q.UpdateMaterialPartial(ctx, database.UpdateMaterialPartialParams{
		Name:        utils.ToSqlNullString(name),
		Uuid:        uuid,
		Url:         utils.ToSqlNullString(url),
		Description: utils.ToSqlNullString(description),
	})
	if err != nil {
		return nil, err
	}
	return &material, nil
}
