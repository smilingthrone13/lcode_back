package webapi

import (
	"lcode/config"
	"lcode/internal/infra/webapi/judge"
)

type (
	InitParams struct {
		Config *config.Config
	}

	APIs struct {
		Judge *judge.API
	}
)

func New(p *InitParams) *APIs {
	judgeAPI := judge.New(&p.Config.JudgeConfig)

	return &APIs{
		Judge: judgeAPI,
	}
}
