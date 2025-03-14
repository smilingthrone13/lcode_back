package article

import (
	"errors"
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	accessMiddleware "lcode/internal/handler/middleware/access"
	"lcode/internal/handler/middleware/article"
	authMiddleware "lcode/internal/handler/middleware/auth"
	articleService "lcode/internal/service/article"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"lcode/pkg/struct_errors"
	"log/slog"
	"net/http"
)

type (
	Middlewares struct {
		Auth    *authMiddleware.Middleware
		Access  *accessMiddleware.Middleware
		Article *article.Middleware
	}

	Services struct {
		Article articleService.Article
	}

	Handler struct {
		config   *config.Config
		logger   *slog.Logger
		services *Services
	}
)

func New(cfg *config.Config, logger *slog.Logger, services *Services) *Handler {
	return &Handler{
		config:   cfg,
		logger:   logger,
		services: services,
	}
}

func (h *Handler) Register(middlewares *Middlewares, httpServer *gin.Engine) {
	practiceGroup := httpServer.Group("/articles/practice")
	{
		practiceGroup.GET("", h.getPracticeArticle)

		practiceGroup.PATCH("", middlewares.Article.ValidateUpdatePracticeArticleInput, h.updateArticle)

		practiceGroup.DELETE("", h.deletePracticeArticle)
	}

	articleGroup := httpServer.Group("/articles", middlewares.Access.UserIdentity)
	{
		articleGroup.GET("/available_attributes", h.getAvailableArticleAttributes)

		articleGroup.GET("/list", middlewares.Article.ValidateArticleListByParamsInput, h.getArticleList)

		articleGroup.GET("/:id", middlewares.Article.ValidateArticleGetInput, h.getArticle)

		adminArticleGroup := articleGroup.Group("", middlewares.Auth.CheckAdminAccess)
		{
			adminArticleGroup.POST("/", middlewares.Article.ValidateCreateArticleInput, h.createArticle)

			adminArticleGroup.PATCH("/:id", middlewares.Article.ValidateUpdateArticleInput, h.updateArticle)

			adminArticleGroup.DELETE("/:id", middlewares.Article.ValidateDeleteArticleInput, h.deleteArticle)
		}
	}
}

func (h *Handler) createArticle(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.ArticleCreateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	a, err := h.services.Article.Create(c.Request.Context(), dto.Input)
	if err != nil {
		var errExist *struct_errors.ErrExist
		if errors.As(err, &errExist) {
			http_helper.NewErrorResponse(c, http.StatusConflict, errExist.Msg)

			return
		}

		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusCreated, a)
}

func (h *Handler) updateArticle(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.ArticleUpdateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	a, err := h.services.Article.Update(c.Request.Context(), dto.Input)
	if err != nil {
		var errExist *struct_errors.ErrExist
		if errors.As(err, &errExist) {
			http_helper.NewErrorResponse(c, http.StatusConflict, errExist.Msg)

			return
		}

		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, a)
}

func (h *Handler) deleteArticle(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.ArticleDeleteDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.Article.Delete(c.Request.Context(), dto.ID)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, map[string]string{"message": "Successful operation"})
}

func (h *Handler) getArticle(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.ArticleGetDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	a, err := h.services.Article.GetByID(c.Request.Context(), dto.ID)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, a)
}

func (h *Handler) getArticleList(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.ArticleParamsDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	articles, err := h.services.Article.GetAllByParams(c.Request.Context(), dto.Input)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, articles)
}

func (h *Handler) getAvailableArticleAttributes(c *gin.Context) {
	aa, err := h.services.Article.GetAvailableAttributes(c.Request.Context())
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, aa)
}

// Practice Interactions

func (h *Handler) getPracticeArticle(c *gin.Context) {
	a, err := h.services.Article.GetByID(c.Request.Context(), domain.PracticeArticleID)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, a)
}

func (h *Handler) deletePracticeArticle(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"message": "Successful operation"})
}
