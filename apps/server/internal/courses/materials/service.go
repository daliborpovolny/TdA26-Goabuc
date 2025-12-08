package materials

import (
	"context"
	"fmt"
	"strings"
	database "tourbackend/internal/database/gen"
)

type Service struct {
	q          *database.Queries
	staticPath string
}

func NewService(queries *database.Queries, staticPath string) *Service {
	return &Service{queries, staticPath}
}

func (s *Service) ListMaterials(courseId string, ctx context.Context) ([]Material, error) {

	dbMaterials, err := s.q.ListAllMaterialsOfCourse(ctx, courseId)
	if err != nil {
		return nil, FailedToFetchMaterials
	}

	formattedMaterials := make([]Material, 0, len(dbMaterials))
	for _, material := range dbMaterials {
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
	return formattedMaterials, nil
}
