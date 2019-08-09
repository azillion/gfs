package gfs

import (
	"github.com/sirupsen/logrus"
)

// Service holds the repository and params of the service
type Service struct {
	repository Repository
	params     *Params
}

// GetFiles from NOMADS
func (s *Service) GetFiles() error {
	baseURL, err := s.repository.GetBaseURL()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Debug(baseURL)

	_, err = s.repository.GetURIs()
	if err != nil {
		logrus.Fatal(err)
	}

	return nil
}

// NewService creates a new gfs service
func NewService(p *Params) *Service {
	r := NewRepository(p.RepositoryType)
	if r == nil {
		logrus.Fatal("no repository of that type")
	}
	err := r.LoadParams(p)
	if err != nil {
		logrus.Fatalf("error loading params: %v", err)
	}
	return &Service{
		repository: r,
		params:     p,
	}
}
