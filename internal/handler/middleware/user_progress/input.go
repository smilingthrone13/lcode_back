package user_progress

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

func (m *Middleware) ValidateGetStatisticsInput(c *gin.Context) {
	qType := c.Query("type")

	if qType == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Type is required")

		return
	}

	if qType != domain.StatisticCategory && qType != domain.StatisticDifficulty {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Invalid type")

		return
	}

	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	dto := domain.GetStatisticsDTO{
		UserID: user.ID,
		Type:   qType,
	}

	c.Set(domain.DtoCtxKey, dto)
}
