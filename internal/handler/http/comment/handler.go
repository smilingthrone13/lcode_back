package comment

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	accessMiddleware "lcode/internal/handler/middleware/access"
	commentMiddleware "lcode/internal/handler/middleware/comment"
	commentService "lcode/internal/service/comment"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
)

type (
	Middlewares struct {
		Access  *accessMiddleware.Middleware
		Comment *commentMiddleware.Middleware
	}

	Services struct {
		Comment commentService.Comment
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
	commentGroup := httpServer.Group("/comments", middlewares.Access.UserIdentity)
	{
		commentGroup.GET("/:entity_id", middlewares.Comment.ValidateThreadsListByParamsInput, h.getThreadList)

		commentGroup.POST("/", middlewares.Comment.ValidateCreateCommentInput, h.createComment)

		commentGroup.PATCH("/:comment_id", middlewares.Comment.ValidateUpdateCommentInput, h.updateComment)

		commentGroup.DELETE("/:comment_id", middlewares.Comment.ValidateDeleteCommentInput, h.deleteComment)
	}
}

func (h *Handler) createComment(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.CommentCreateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	comm, err := h.services.Comment.Create(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusCreated, comm)
}

func (h *Handler) updateComment(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.CommentUpdateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	comm, err := h.services.Comment.Update(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, comm)
}

func (h *Handler) deleteComment(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.CommentDeleteDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.Comment.Delete(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, map[string]string{"message": "Successful operation"})
}

func (h *Handler) getThreadList(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.CommentParamsDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	threads, err := h.services.Comment.GetThreadsByParamsAndEntityID(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, threads)
}
