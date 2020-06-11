package docApi

import (
	"context"
	"io/ioutil"
	"log"
	"os"
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
	exists := true

	routePath := path.Join(doc.Squad, doc.Projeto, doc.Versao)
	folderPath := path.Join(FilesFolder, routePath)
	filePath := path.Join(folderPath, FileName)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			exists = false
		} else {
			log.Printf("File %s stat error: %v", filePath, err)
		}
	}

	doc.Doc = []byte(path.Join(routePath, FileName))

	if exists {
		return doc, nil
	}

	res, err := s.repo.Find(ctx, doc.Squad, doc.Projeto, doc.Versao)
	if err != nil {
		return nil, errors.Wrap(err, "could not find doc")
	}

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "could not create the folderPath")
	}

	err = ioutil.WriteFile(filePath, res.Doc, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "could not write on file")
	}

	return doc, nil
}

func (s *Service) FindBySquad(ctx context.Context, squad string, limit string, offset string) (*SearchResult, error) {
	return s.repo.FindBySquad(ctx, squad, limit, offset)
}

func (s *Service) Delete(ctx context.Context, doc *DocApi) error {
	filePath := path.Join(FilesFolder, doc.Squad, doc.Projeto, doc.Versao, FileName)

	err := os.Remove(filePath)
	if err != nil {
		log.Printf("could not delete %s: %v", filePath, err)
	}

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
