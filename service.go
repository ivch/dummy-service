package dummy_service

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ivch/dummy-service/repository"
)

type Service interface {
	Create(req []byte) ([]byte, error)
}

type service struct {
	repo Repository
}

type Repository interface {
	Insert(e *repository.Event) error
}

func New(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(req []byte) ([]byte, error) {
	var e repository.Event

	if err := json.Unmarshal(req, &e); err != nil {
		return nil, err
	}

	if err := s.repo.Insert(&e); err != nil {
		return nil, errors.Wrap(err, "failed inserting event")
	}

	return json.Marshal(e)
}
