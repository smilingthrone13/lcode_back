package problem

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	accessMiddleware "lcode/internal/handler/middleware/access"
	authMiddleware "lcode/internal/handler/middleware/auth"
	problemMiddleware "lcode/internal/handler/middleware/problem"
	"lcode/internal/manager/problem_manager"
	"log/slog"
)

type (
	Middlewares struct {
		Access  *accessMiddleware.Middleware
		Auth    *authMiddleware.Middleware
		Problem *problemMiddleware.Middleware
	}

	Managers struct {
		Problem *problem_manager.Manager
	}

	Handler struct {
		config   *config.Config
		logger   *slog.Logger
		managers *Managers
	}
)

func New(cfg *config.Config, logger *slog.Logger, managers *Managers) *Handler {
	return &Handler{
		config:   cfg,
		logger:   logger,
		managers: managers,
	}
}

func (h *Handler) Register(middlewares *Middlewares, httpServer *gin.Engine) {
}
