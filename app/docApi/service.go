package docApi

import (
	"context"

	"github.com/pkg/errors"
)

type Service struct {
	repo Repository
}

func (s *Service) Create(ctx context.Context, doc *DocApi) (*DocApi, error) {
	newDoc, err := s.repo.Create(ctx, doc)
	if err != nil {
		return &DocApi{}, errors.Wrap(err, "could not create a document")
	}

	return newDoc, nil
}

func (s *Service) Find(ctx context.Context, squad string, projeto string, versao string) (*DocApi, error) {
	return s.repo.Find(ctx, squad, projeto, versao)
}

func (s *Service) Delete(ctx context.Context, squad string, projeto string, versao string) error {
	return s.repo.Delete(ctx, squad, projeto, versao)
}

func (s *Service) FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func NewService(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}
