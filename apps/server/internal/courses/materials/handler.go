package materials

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	db "tourbackend/internal/database/gen"
	"tourbackend/internal/handlers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

type MaterialsHandlers struct {
	*handlers.Handler
	staticPath string
}

func NewMaterialsHandlers(staticPath string, queries *db.Queries, isDeployed bool) *MaterialsHandlers {
	return &MaterialsHandlers{
		handlers.NewHandler(queries, isDeployed),
		staticPath,
	}
}

type Material interface{}

type FileMaterial struct {
	Uuid        string
	Type        string
	Name        string
	Description string
	FileUrl     string
	MimeType    string
	SizeBytes   int
}

type UrlMaterial struct {
	Uuid        string
	Type        string
	Name        string
	Description string
	Url         string
	FaviconUrl  string
}

func (h *MaterialsHandlers) ListMaterials(c echo.Context) error {
	r := h.NewReqCtx(c)
	courseId := c.Param("courseId")

	materials, err := r.Queries.ListAllMaterialsOfCourse(r.Ctx, courseId)
	if err != nil {
		return r.Error(http.StatusInternalServerError, "Failed to fetch courses from the database.")
	}

	formattedMaterials := make([]Material, 0, len(materials))
	for _, material := range materials {
		if strings.HasPrefix(material.Url, "/this-site") {
			formattedMaterials = append(formattedMaterials, FileMaterial{
				Uuid:        material.Uuid,
				Type:        "url",
				Description: material.Description,
				FileUrl:     material.Url,
				MimeType:    "mime",
				SizeBytes:   0,
			})
		} else {

			urlParts := strings.SplitAfterN(material.Url, "/", 4)
			fmt.Println(urlParts)

			urlBaseSite := urlParts[0] + urlParts[1] + urlParts[2]
			fmt.Println(urlBaseSite)

			faviconUrl := urlBaseSite + "/favicon.ico"
			fmt.Println(faviconUrl)

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

	return c.String(http.StatusOK, "List materials, Coures: "+courseId)
}

func (h *MaterialsHandlers) CreateMaterial(c echo.Context) error {
	r := h.NewReqCtx(c)

	courseId := c.Param("courseId")

	materialType := c.FormValue("type")
	name := c.FormValue("name")
	description := c.FormValue("description")

	if materialType == "" || (materialType != "url" && materialType != "file") {
		return r.Error(http.StatusBadRequest, "type of the material must be provided, either 'file' or 'url'")
	}

	if name == "" {
		return r.Error(http.StatusBadRequest, "name of the material must be provided")
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return r.Error(http.StatusBadRequest, "file is required")
	}

	if fileHeader.Size > MAX_SIZE {
		return r.Error(http.StatusBadRequest, "file is too big, max 30MB")
	}

	src, err := fileHeader.Open()
	if err != nil {
		fmt.Println("error when opening file header", err)
		return r.Error(http.StatusInternalServerError, "failed to open the file header")
	}
	defer src.Close()

	// read the mimetype
	buf := make([]byte, 512)
	n, _ := src.Read(buf)

	mimeType := http.DetectContentType(buf[:n])
	fmt.Println(mimeType)

	if !ALLOWED_FILES[mimeType] {
		fmt.Println("forbiden file type: ", mimeType)
		return r.Error(http.StatusBadRequest, "forbiden file type")
	}

	// 	rewind the cursor
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return err
	}

	pathToFolder := h.staticPath + "/uploads/" + courseId + "/materials"
	materialId := uuid.NewString()
	pathToFile := pathToFolder + "/" + materialId

	err = os.Mkdir(pathToFolder, 0755)
	if err != nil {
		fmt.Println("wrong path to folder")
		panic(err)
	}
	dst, err := os.Create(pathToFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	req := c.Request()
	scheme := c.Scheme()
	host := req.Host

	_, err = r.Queries.CreateMaterial(r.Ctx, db.CreateMaterialParams{
		Uuid:        materialId,
		Name:        name,
		Description: description,
		Url:         scheme + "://" + host + "/static/uploads/" + courseId + "/materials/" + materialId,
	})
	if err != nil {
		fmt.Println("failed to create material in the db", err)
		return r.Error(http.StatusInternalServerError, "failed to save the material to the db")
	}

	return c.String(http.StatusOK, "Create material, Coures: "+courseId)
}

func (h *MaterialsHandlers) UpdateMaterial(c echo.Context) error {
	courseId := c.Param("courseId")
	materialId := c.Param("materialId")

	return c.String(http.StatusOK, "Update material, Course: "+courseId+" Material:"+materialId)
}

func (h *MaterialsHandlers) DeleteMaterial(c echo.Context) error {
	courseId := c.Param("courseId")
	materialId := c.Param("materialId")

	return c.String(http.StatusOK, "Delete material: Course:"+courseId+" Material:"+materialId)
}
