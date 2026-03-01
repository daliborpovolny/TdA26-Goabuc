package headings

import (
	"context"

	db "tourbackend/internal/database/gen"
	"tourbackend/internal/feeds"
)

type Service struct {
	q            *db.Queries
	staticPath   string
	feedsService *feeds.Service
}

func NewService(queries *db.Queries, staticPath string, feedsService *feeds.Service) *Service {
	return &Service{queries, staticPath, feedsService}
}

func (s *Service) AssignHeadingToModule(headingId string, moduleId string, order int, ctx context.Context) (db.HeadingToModule, error) {

	hm, err := s.q.AssignHeadingToModule(ctx, db.AssignHeadingToModuleParams{
		ModuleUuid:  moduleId,
		HeadingUuid: headingId,
		Order:       int64(order),
	})
	if err != nil {
		return db.HeadingToModule{}, err
	}

	return hm, nil
}

func (s *Service) ChangeHeadingInModuleOrder(headingId string, moduleId string, order int, ctx context.Context) (db.HeadingToModule, error) {

	hm, err := s.q.ChangeHeadingInModuleOrder(ctx, db.ChangeHeadingInModuleOrderParams{
		ModuleUuid:  moduleId,
		HeadingUuid: headingId,
		Order:       int64(order),
	})
	if err != nil {
		return db.HeadingToModule{}, err
	}

	return hm, nil
}

func (s *Service) RemoveHeadingToModule(headingId string, moduleId string, order int, ctx context.Context) error {

	err := s.q.RemoveHeadingFromModule(ctx, db.RemoveHeadingFromModuleParams{
		ModuleUuid:  moduleId,
		HeadingUuid: headingId,
	})
	if err != nil {
		return err
	}
	return nil
}
