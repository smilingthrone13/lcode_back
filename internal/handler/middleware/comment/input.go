package comment

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
)

const (
	articleOriginInput = "article"
	problemOriginInput = "problem"
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

func (m *Middleware) ValidateCreateCommentInput(c *gin.Context) {
	var dto domain.CommentCreateDTO

	origin := c.Param("origin_type")
	switch origin {
	case articleOriginInput:
		dto.OriginType = domain.ArticleOriginType
	case problemOriginInput:
		dto.OriginType = domain.TaskOriginType
	default:
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Unknown origin type")

		return
	}

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if dto.Input.EntityID == "" || dto.Input.Text == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Entity ID and comment text are required")

		return
	}

	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	dto.Input.AuthorID = user.ID

	c.Set(domain.DtoCtxKey, dto)

}

func (m *Middleware) ValidateUpdateCommentInput(c *gin.Context) {
	var dto domain.CommentUpdateDTO

	origin := c.Param("origin_type")
	switch origin {
	case articleOriginInput:
		dto.OriginType = domain.ArticleOriginType
	case problemOriginInput:
		dto.OriginType = domain.TaskOriginType
	default:
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Unknown origin type")

		return
	}

	dto.Input.ID = c.Param("comment_id")
	if dto.Input.ID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Comment ID is required")

		return
	}

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if dto.Input.Text == nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Comment text is required")

		return
	}

	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	dto.User = user

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateDeleteCommentInput(c *gin.Context) {
	var dto domain.CommentDeleteDTO

	origin := c.Param("origin_type")
	switch origin {
	case articleOriginInput:
		dto.OriginType = domain.ArticleOriginType
	case problemOriginInput:
		dto.OriginType = domain.TaskOriginType
	default:
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Unknown origin type")

		return
	}

	dto.ID = c.Param("comment_id")
	if dto.ID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Comment ID is required")

		return
	}

	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	dto.User = user

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateThreadsListByParamsInput(c *gin.Context) {
	var dto domain.CommentParamsDTO

	origin := c.Param("origin_type")
	switch origin {
	case articleOriginInput:
		dto.OriginType = domain.ArticleOriginType
	case problemOriginInput:
		dto.OriginType = domain.TaskOriginType
	default:
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Unknown origin type")

		return
	}

	dto.EntityID = c.Param("entity_id")
	if dto.EntityID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Entity ID is required")

		return
	}

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}
