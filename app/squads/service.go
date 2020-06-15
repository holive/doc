package squads

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Service struct {
	repo Repository
}

func (s *Service) Create(ctx context.Context, name string) (Squad, error) {
	squadKey := s.createHash("parrot não é papagaio")

	newSquad := Squad{
		Name: strings.Trim(name, " "),
		Key:  squadKey,
	}

	return s.repo.Create(ctx, newSquad)
}

func (s *Service) VerifyUserKey(ctx context.Context, squad Squad) (bool, error) {
	sqd, err := s.repo.GetByKey(ctx, squad.Key)
	if err != nil {
		return false, errors.Wrap(err, "could not get squad using the key")
	}

	if sqd.Name == squad.Name {
		return true, nil
	}

	return false, nil
}

func (s *Service) createHash(key string) string {
	tu := time.Now().UnixNano()
	str := strconv.FormatInt(tu, 10)

	h := sha1.New()
	h.Write([]byte(str + key))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func NewService(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}
