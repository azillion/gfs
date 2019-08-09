package gfs

import (
	"fmt"
)

// Service holds the repository and params of the service
type Service struct {
	repository Repository
	params     *Params
}

// GetFiles get files from NOMADS
func (s *Service) GetFiles() error {
	if s.params.TimeFrame == AllTimeFrames {
		_, err := s.repository.GetURIs()
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// NewService creates a new gfs service
func NewService(p *Params) *Service {
	r := NewRepository(p.RepositoryType)
	if r == nil {
		panic(fmt.Errorf("no repository of that type"))
	}
	err := r.LoadParams(p)
	if err != nil {
		panic(fmt.Errorf("error loading params: %v", err))
	}
	return &Service{
		repository: r,
		params:     p,
	}
}
