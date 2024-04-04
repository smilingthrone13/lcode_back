package general

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	generalMiddleware "lcode/internal/handler/middleware/general"
	"log/slog"
)

type (
	Middlewares struct {
		General *generalMiddleware.Middleware
	}

	Services struct {
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
}
