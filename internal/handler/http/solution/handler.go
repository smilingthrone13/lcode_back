package solution

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	accessMiddleware "lcode/internal/handler/middleware/access"
	solutionMiddleware "lcode/internal/handler/middleware/solution"
	"lcode/internal/manager/solution_manager"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
)

type (
	Middlewares struct {
		Access   *accessMiddleware.Middleware
		Solution *solutionMiddleware.Middleware
	}

	Services struct {
		SolutionManager solution_manager.SolutionManager
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
	solutionGroup := httpServer.Group("/solution")
	{
		solutionGroup.POST(
			"/",
			middlewares.Access.UserIdentity,
			middlewares.Solution.ValidateCreateSolutionInput,
			h.createSolution,
		)
	}
}

func (h *Handler) createSolution(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.CreateSolutionDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	sol, err := h.services.SolutionManager.CreateSolution(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, sol)
}
