package user_progress

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	accessMiddleware "lcode/internal/handler/middleware/access"
	userProgressMiddleware "lcode/internal/handler/middleware/user_progress"
	userProgress "lcode/internal/service/user_progress"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
)

type (
	Middlewares struct {
		Access       *accessMiddleware.Middleware
		UserProgress *userProgressMiddleware.Middleware
	}

	Services struct {
		UserProgress userProgress.UserProgress
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
	userProgressGroup := httpServer.Group("/progress", middlewares.Access.UserIdentity)
	{
		userProgressGroup.GET(
			"/",
			h.getUserProgress,
		)

		userProgressGroup.GET(
			"/statistics",
			middlewares.UserProgress.ValidateGetStatisticsInput,
			h.getUserStatistics,
		)
	}
}

func (h *Handler) getUserProgress(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	up, err := h.services.UserProgress.GetProgressByUserID(c.Request.Context(), dto.ID)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, up)
}

func (h *Handler) getUserStatistics(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.GetStatisticsDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	us, err := h.services.UserProgress.GetStatisticsByUserID(c.Request.Context(), dto.UserID, dto.Type)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, us)
}
