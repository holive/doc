package docApi

import (
	"context"
	"io/ioutil"
	"path"

	"github.com/pkg/errors"
)

type Service struct {
	repo Repository
}

func (s *Service) Create(ctx context.Context, folderPath string, filename string, document *DocApi) error {
	doc, err := ioutil.ReadFile(path.Join(folderPath, filename))
	if err != nil {
		return errors.Wrap(err, "could not open document file")
	}

	document.Doc = doc

	err = s.repo.Create(ctx, document)
	if err != nil {
		return errors.Wrap(err, "could not create a document")
	}

	return nil
}

func (s *Service) Find(ctx context.Context, doc *DocApi) (*DocApi, error) {
	return s.repo.Find(ctx, doc.Squad, doc.Projeto, doc.Versao)
}

func (s *Service) Delete(ctx context.Context, doc *DocApi) error {
	return s.repo.Delete(ctx, doc.Squad, doc.Projeto, doc.Versao)
}

func (s *Service) FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func NewService(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}
