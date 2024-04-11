package solution

import (
	"context"
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/internal/service/solution"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
	"slices"
)

type (
	Solution interface {
		SolutionByID(ctx context.Context, id string) (sol domain.Solution, err error)
	}
)

type (
	Services struct {
		Solution solution.Solution
	}

	Middleware struct {
		cfg      *config.Config
		logger   *slog.Logger
		services *Services
	}
)

func New(cfg *config.Config, logger *slog.Logger, services *Services) *Middleware {
	return &Middleware{
		cfg:      cfg,
		logger:   logger,
		services: services,
	}
}

type createSolutionInput struct {
	TaskID     string              `json:"task_id"`
	LanguageID domain.LanguageType `json:"language_id"`
	Code       string              `json:"code"`
}

func (m *Middleware) ValidateCreateSolutionInput(c *gin.Context) {
	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	var inp createSolutionInput

	err = c.ShouldBindJSON(&inp)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if !slices.Contains(domain.AvailableLanguageIds, inp.LanguageID) {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "not supported language")
	}

	dto := domain.CreateSolutionDTO{
		TaskID:     inp.TaskID,
		LanguageID: inp.LanguageID,
		Code:       inp.Code,
		User:       user,
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateGetSolutionsInput(c *gin.Context) {
	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	dto := domain.GetSolutionsDTO{
		TaskID: c.Param("task_id"),
		User:   user,
	}
	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateGetSolutionResultsInput(c *gin.Context) {
	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	dto := domain.GetSolutionResultsDTO{
		SolutionID: c.Param("id"),
		User:       user,
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateGetSolutionCodeInput(c *gin.Context) {
	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	dto := domain.GetSolutionCodeDTO{
		SolutionID: c.Param("id"),
		User:       user,
	}

	c.Set(domain.DtoCtxKey, dto)
}
