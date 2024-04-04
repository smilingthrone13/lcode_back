package repository

import (
	"lcode/config"
	"lcode/internal/infra/repository/general"
	"lcode/pkg/postgres"
)

type (
	InitParams struct {
		Config *config.Config
		DB     *postgres.DbManager
	}

	Repositories struct {
		General *general.Repository
	}
)

func New(p *InitParams) *Repositories {
	return &Repositories{
		General: general.New(p.DB),
	}
}
