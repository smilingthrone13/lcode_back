package article

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
)

type (
	Middleware struct {
		cfg    *config.Config
		logger *slog.Logger
	}
)

func New(cfg *config.Config, logger *slog.Logger) *Middleware {
	return &Middleware{
		cfg:    cfg,
		logger: logger,
	}
}

func (m *Middleware) ValidateCreateArticleInput(c *gin.Context) {
	var dto domain.ArticleCreateDTO

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if dto.Input.Title == "" || dto.Input.Content == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Title and content are required")

		return
	}

	if dto.Input.Categories == nil {
		dto.Input.Categories = []string{}
	}

	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	dto.Input.AuthorID = user.ID

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateUpdateArticleInput(c *gin.Context) {
	var dto domain.ArticleUpdateDTO

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if dto.Input.Title == nil && dto.Input.Content == nil && dto.Input.Categories == nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "No update data provided")

		return
	}

	dto.Input.ID = c.Param("id")
	if dto.Input.ID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Article ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateDeleteArticleInput(c *gin.Context) {
	var dto domain.ArticleDeleteDTO

	dto.ID = c.Param("id")
	if dto.ID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Article ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateArticleListByParamsInput(c *gin.Context) {
	var inp domain.ArticleParamsInput

	if err := c.ShouldBindJSON(&inp); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	categories, ok := c.GetQueryArray("category")
	if !ok {
		categories = []string{}
	}

	filter := domain.ArticleFilter{
		Search:     c.Query("search"),
		Categories: categories,
	}

	data := domain.ArticleParams{
		Filter:     filter,
		Sort:       inp.Sort,
		Pagination: inp.Pagination,
	}

	c.Set(domain.DtoCtxKey, domain.ArticleParamsDTO{Input: data})
}

func (m *Middleware) ValidateArticleGetInput(c *gin.Context) {
	dto := domain.ArticleGetDTO{
		ID: c.Param("id"),
	}

	if dto.ID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Article ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

// Practice

func (m *Middleware) ValidateUpdatePracticeArticleInput(c *gin.Context) {
	var dto domain.ArticleUpdateDTO

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if dto.Input.Title == nil && dto.Input.Content == nil && dto.Input.Categories == nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "No update data provided")

		return
	}

	dto.Input.ID = domain.PracticeArticleID

	c.Set(domain.DtoCtxKey, dto)
}
