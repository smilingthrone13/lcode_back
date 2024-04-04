package authorization

import (
	"lcode/config"
	"lcode/internal/domain"
	"sync"
)

type (
	Cache struct {
		mu           sync.RWMutex
		users        []domain.User
		userIdToUser map[string]domain.User
	}

	Service struct {
		config *config.Config
		cache  *Cache
	}
)

func NewService(conf *config.Config) *Service {
	cache := &Cache{userIdToUser: make(map[string]domain.User, 1000)}

	// todo: get users from db into cache

	s := &Service{
		config: conf,
		cache:  cache,
	}

	return s
}

func (s *Service) UserByID(id string) (user domain.User, ok bool) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Users() ([]domain.User, error) {
	//TODO implement me
	panic("implement me")
}
