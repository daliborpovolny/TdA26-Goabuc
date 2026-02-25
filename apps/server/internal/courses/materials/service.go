package materials

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	db "tourbackend/internal/database/gen"
	"tourbackend/internal/feeds"
	"tourbackend/internal/utils"

	"github.com/gabriel-vasile/mimetype"
)

// this variable controls whether a material can exist withouth being part of a module
var MATERIAL_CAN_EXIST_ALONE = true

// max file size in bytes
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
	"image/jpeg": ".jpeg",
	"image/gif":  ".gif",
	"video/mp4":  ".mp4",
	"audio/mpeg": ".mp3",
}

var EXT_TO_MIME = map[string]string{
	".pdf":  "application/pdf",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".txt":  "text/plain",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	"gif":   "image/gif",
	".mp4":  "video/mp4",
	".mp3":  "audio/mpeg",
}

type Material interface {
	GetType() string
	GetUuid() string
}

type FileMaterial struct {
	Uuid        string `json:"uuid"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`

	FileUrl   string `json:"fileUrl"`
	MimeType  string `json:"mimeType"`
	SizeBytes int    `json:"sizeBytes"`
}

func (f FileMaterial) GetType() string {
	return f.Type
}

func (f FileMaterial) GetUuid() string {
	return f.Uuid
}

type UrlMaterial struct {
	Uuid        string `json:"uuid"`
	ModuleUuid  string `json:"module_uuid"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Url        string `json:"url"`
	FaviconUrl string `json:"faviconUrl"`
}

func (f UrlMaterial) GetType() string {
	return f.Type
}

func (f UrlMaterial) GetUuid() string {
	return f.Uuid
}

type Service struct {
	q            *db.Queries
	staticPath   string
	feedsService *feeds.Service
}

func NewService(queries *db.Queries, staticPath string, feedsService *feeds.Service) *Service {
	return &Service{queries, staticPath, feedsService}
}

func (s *Service) deriveFaviconUrl(url string) string {

	if strings.HasPrefix(url, "http") {
		parts := strings.Split(url, "//")

		host := parts[1]

		base := strings.Split(host, "/")[0]

		faviconUrl := base + "/favicon.ico"
		return faviconUrl
	}

	urlParts := strings.Split(url, "/")
	faviconUrl := urlParts[0] + "/favicon.ico"
	return faviconUrl
}

func (s *Service) CheckMaterialExists(materialId string, ctx context.Context) bool {
	_, err := s.q.GetMaterial(ctx, materialId)
	if err != nil {
		if utils.IsNoRowsError(err) {
			return false
		}
		fmt.Println("Check if material exists failed, uuid:", materialId)
		return false
	}
	return true
}

func (s *Service) ListMaterials(courseId string, host string, scheme string, ctx context.Context) ([]Material, error) {

	dbMaterials, err := s.q.ListAllMaterialsOfCourse(ctx, courseId)
	if err != nil {
		return nil, err
	}

	if len(dbMaterials) == 0 {
		ok, err := s.q.CheckCourseExists(ctx, courseId)
		if err != nil {
			return nil, err
		}

		if ok != 1 {
			return nil, ErrCourseNotFound
		}
		return []Material{}, nil
	}

	formattedMaterials := make([]Material, 0, len(dbMaterials))
	for _, material := range dbMaterials {

		if material.Type == "file" {

			formattedMaterials = append(formattedMaterials, FileMaterial{
				Uuid:        material.Uuid,
				Type:        "file",
				Name:        material.Name,
				Description: material.Description,
				FileUrl:     material.Url,
				MimeType:    material.MimeType.String,
				SizeBytes:   int(material.ByteSize.Int64),
			})
		} else {

			formattedMaterials = append(formattedMaterials, UrlMaterial{
				Uuid:        material.Uuid,
				Type:        "url",
				Name:        material.Name,
				Description: material.Description,
				Url:         material.Url,
				FaviconUrl:  material.FaviconUrl.String,
			})
		}
	}
	return formattedMaterials, nil
}

func (s *Service) getMimeType(src multipart.File) (string, error) {

	m, err := mimetype.DetectReader(src)
	if err != nil {
		return "", err
	}

	mimeType := m.String()

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", nil
	}

	return mimeType, nil
}

// checks if the file type is allowed, returns bool, mime and error
func (s *Service) isFileAllowed(src multipart.File) (bool, string, error) {

	mime, err := s.getMimeType(src)
	if err != nil {
		return false, "", err
	}

	for key := range ALLOWED_FILES {
		if strings.Contains(mime, key) {
			return true, key, nil
		}
	}
	return false, "", nil
}

func (s *Service) isFileAllowedHeader(file *multipart.FileHeader) (string, bool) {
	if file == nil {
		return "", false
	}

	filename := file.Filename
	ext := strings.ToLower(filepath.Ext(filename))
	mime, ok := EXT_TO_MIME[ext]
	if ok {
		return mime, ok
	}

	mime = file.Header.Get("Content-Type")
	ok = ALLOWED_FILES[mime]
	if ok {
		return mime, true
	}
	return "", false
}

// removes old material file when a new one is uploaded
func (s *Service) removeOldMaterialFile(materialId string, courseId string) error {

	pathToMaterialsFolder := s.staticPath + "/uploads/" + courseId + "/materials"

	entries, err := os.ReadDir(pathToMaterialsFolder)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasPrefix(name, materialId) {
			err := os.Remove(filepath.Join(pathToMaterialsFolder, name))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// saves given file to the static folder
func (s *Service) saveFileToStatic(
	materialId string,
	courseId string,
	src multipart.File,
	ext string,
) error {

	pathToUploads := s.staticPath + "/uploads"
	pathToCourseFolder := pathToUploads + "/" + courseId
	pathToMaterialsFolder := pathToCourseFolder + "/materials"

	pathToFile := pathToMaterialsFolder + "/" + materialId + ext

	err := os.Mkdir(pathToUploads, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			fmt.Println("failed to create materials folder, wrong path", err)
			return err
		}
	}

	err = os.Mkdir(pathToCourseFolder, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			fmt.Println("failed to create course folder, wrong path", err)
			return err
		}
	}

	err = os.Mkdir(pathToMaterialsFolder, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			fmt.Println("failed to create materials folder, wrong path", err)
			return err
		}
	}

	dst, err := os.Create(pathToFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateFileMaterial(req *CreateFileMaterialRequest, materialId string, fileHeader *multipart.FileHeader, scheme string, host string, ctx context.Context) (Material, error) {

	if fileHeader.Size > MAX_SIZE {
		return nil, ErrFileTooBig
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	// ok, mime, err := s.isFileAllowed(src)
	// if err != nil {
	// 	return nil, err
	// }
	mime, ok := s.isFileAllowedHeader(fileHeader)

	if !ok {
		return nil, ErrFileTypeForbidden
	}

	fileExt := MIME_TO_EXT[mime]

	err = s.saveFileToStatic(materialId, req.CourseId, src, fileExt)

	url := scheme + "://" + host + "/api/static/uploads/" + req.CourseId + "/materials/" + materialId + fileExt

	now := time.Now().Unix()

	dbMat, err := s.q.CreateMaterial(ctx, db.CreateMaterialParams{
		Uuid:        materialId,
		CourseUuid:  req.CourseId,
		Name:        req.Name,
		Description: req.Description,
		Url:         url,
		Type:        "file",
		MimeType:    sql.NullString{String: mime, Valid: true},
		ByteSize:    sql.NullInt64{Int64: fileHeader.Size, Valid: true},
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, err
	}

	s.feedsService.CreateAutomaticPost("New file material: "+req.Name+" published", req.CourseId, ctx)
	return FileMaterial{
		Uuid:        dbMat.Uuid,
		Type:        dbMat.Type,
		Name:        dbMat.Name,
		Description: dbMat.Description,
		FileUrl:     dbMat.Url,
		MimeType:    dbMat.MimeType.String,
		SizeBytes:   int(dbMat.ByteSize.Int64),
	}, nil
}

type CreateUrlMaterialParams struct {
	uuid        string
	name        string
	description string
	courseId    string

	url string
}

func (s *Service) CreateUrlMaterial(req CreateUrlMaterialRequest, materialId string, ctx context.Context) (Material, error) {

	now := time.Now().Unix()

	dbMat, err := s.q.CreateMaterial(ctx, db.CreateMaterialParams{
		Uuid:        materialId,
		CourseUuid:  req.CourseId,
		Name:        req.Name,
		Description: req.Description,
		Url:         req.Url,
		Type:        "url",
		FaviconUrl:  sql.NullString{String: s.deriveFaviconUrl(req.Url), Valid: true},
		UpdatedAt:   now,
		CreatedAt:   now,
	})
	if err != nil {
		return nil, err
	}

	s.feedsService.CreateAutomaticPost("New url material: "+req.Name+" published", req.CourseId, ctx)
	return UrlMaterial{
		Uuid:        dbMat.Uuid,
		Type:        dbMat.Type,
		Name:        dbMat.Name,
		Description: dbMat.Description,
		Url:         dbMat.Url,
		FaviconUrl:  dbMat.FaviconUrl.String,
	}, nil
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

func (s *Service) UpdateFileMaterial(req *UpdateFileMaterialRequest, fileHeader *multipart.FileHeader, scheme string, host string, ctx context.Context) (Material, error) {

	var byteSize = sql.NullInt64{}
	var mimeType *string
	var url *string

	if fileHeader != nil {

		if fileHeader.Size > MAX_SIZE {
			return nil, ErrFileTooBig
		}
		byteSize.Valid = true
		byteSize.Int64 = fileHeader.Size

		src, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}

		// ok, mime, err := s.isFileAllowed(src)
		// if err != nil {
		// 	return nil, err
		// }
		mime, ok := s.isFileAllowedHeader(fileHeader)

		if !ok {
			return nil, ErrFileTypeForbidden
		}

		err = s.removeOldMaterialFile(req.MaterialId, req.CourseId)
		if err != nil {
			return nil, err
		}

		fileExt := MIME_TO_EXT[mime]

		err = s.saveFileToStatic(req.MaterialId, req.CourseId, src, fileExt)
		if err != nil {
			return nil, err
		}
		newUrl := scheme + "://" + host + "/api/static/uploads/" + req.CourseId + "/materials/" + req.MaterialId + fileExt
		url = &newUrl

		mimeType = &mime
	}

	now := time.Now().Unix()

	dbMat, err := s.q.UpdateMaterialPartial(ctx, db.UpdateMaterialPartialParams{
		Uuid:        req.MaterialId,
		Name:        utils.ToSqlNullString(req.Name),
		Description: utils.ToSqlNullString(req.Description),
		Url:         utils.ToSqlNullString(url),
		MimeType:    utils.ToSqlNullString(mimeType),
		ByteSize:    byteSize,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, err
	}

	s.feedsService.CreateAutomaticPost("File material: "+*req.Name+" updated", req.CourseId, ctx)
	return FileMaterial{
		Uuid:        dbMat.Uuid,
		Type:        dbMat.Type,
		Name:        dbMat.Name,
		Description: dbMat.Description,
		FileUrl:     dbMat.Url,
		MimeType:    dbMat.MimeType.String,
		SizeBytes:   int(dbMat.ByteSize.Int64),
	}, nil
}

func (s *Service) UpdateUrlMaterial(req *UpdateUrlMaterialRequest, ctx context.Context) (Material, error) {

	faviconUrl := sql.NullString{}

	if req.Url != nil {
		faviconUrl.Valid = true
		faviconUrl.String = s.deriveFaviconUrl(*req.Url)
	}

	now := time.Now().Unix()

	dbMat, err := s.q.UpdateMaterialPartial(ctx, db.UpdateMaterialPartialParams{
		Uuid:        req.MaterialId,
		Name:        utils.ToSqlNullString(req.Name),
		Description: utils.ToSqlNullString(req.Description),
		Url:         utils.ToSqlNullString(req.Url),
		FaviconUrl:  faviconUrl,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, err
	}

	s.feedsService.CreateAutomaticPost("Url material: "+*req.Name+" updated", req.CourseId, ctx)
	return UrlMaterial{
		Uuid:        dbMat.Uuid,
		Type:        dbMat.Type,
		Name:        dbMat.Name,
		Description: dbMat.Description,
		Url:         dbMat.Url,
		FaviconUrl:  dbMat.FaviconUrl.String,
	}, nil
}

func (s *Service) DeleteMaterial(materialId string, ctx context.Context) error {

	res, err := s.q.DeleteMaterial(ctx, materialId)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return ErrCourseNotFound
	}

	return nil
}

// Material to Module

func (s *Service) AssignMaterialToModule(materialId string, moduleId string, order int, ctx context.Context) (db.MaterialToModule, error) {

	mm, err := s.q.AssignMaterialToModule(ctx, db.AssignMaterialToModuleParams{
		ModuleUuid:   moduleId,
		MaterialUuid: materialId,
		Order:        int64(order),
	})
	if err != nil {
		return db.MaterialToModule{}, err
	}

	return mm, nil
}

func (s *Service) ChangeMaterialInModuleOrder(materialId string, moduleId string, order int, ctx context.Context) (db.MaterialToModule, error) {

	mm, err := s.q.ChangeMaterialInModuleOrder(ctx, db.ChangeMaterialInModuleOrderParams{
		ModuleUuid:   moduleId,
		MaterialUuid: materialId,
		Order:        int64(order),
	})
	if err != nil {
		return db.MaterialToModule{}, err
	}

	return mm, nil
}

func (s *Service) RemoveMaterialToModule(materialId string, moduleId string, order int, ctx context.Context) error {

	err := s.q.RemoveMaterialFromModule(ctx, db.RemoveMaterialFromModuleParams{
		ModuleUuid:   moduleId,
		MaterialUuid: materialId,
	})
	if err != nil {
		return err
	}
	return nil
}
