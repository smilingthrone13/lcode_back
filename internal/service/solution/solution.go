package solution

import (
	"lcode/config"
)

type (
	Service struct {
		config     *config.Config
		repository SolutionRepo
	}
)

func New(conf *config.Config, repository SolutionRepo) *Service {
	return &Service{
		config:     conf,
		repository: repository,
	}
}
