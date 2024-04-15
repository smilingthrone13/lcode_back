package solution

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	accessMiddleware "lcode/internal/handler/middleware/access"
	solutionMiddleware "lcode/internal/handler/middleware/solution"
	"lcode/internal/manager/solution_manager"
	"lcode/internal/service/solution"
	"lcode/internal/service/solution_result"
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
		SolutionManager       solution_manager.SolutionManager
		SolutionService       solution.Solution
		SolutionResultService solution_result.SolutionResult
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
	solutionsGroup := httpServer.Group("/solutions", middlewares.Access.UserIdentity)
	{
		solutionsGroup.GET(
			"/available_statuses",
			h.getAvailableSolutionStatuses,
		)
		solutionsGroup.POST(
			"/",
			middlewares.Solution.ValidateCreateSolutionInput,
			h.createSolution,
		)
		solutionsGroup.GET("/task/:task_id",
			middlewares.Solution.ValidateGetSolutionsInput,
			h.solutions,
		)

		solGroup := solutionsGroup.Group("/:id")
		{
			solGroup.GET(
				"/code",
				middlewares.Solution.ValidateGetSolutionCodeInput,
				middlewares.Solution.CheckSolutionAccess,
				h.solutionCode,
			)

			solGroup.GET(
				"/results",
				middlewares.Access.UserIdentity,
				middlewares.Solution.ValidateGetSolutionResultsInput,
				middlewares.Solution.CheckSolutionAccess,
				h.solutionResults,
			)
		}

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

func (h *Handler) solutions(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.GetSolutionsDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	solutions, err := h.services.SolutionService.SolutionsByUserAndTask(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, solutions)
}

func (h *Handler) solutionResults(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.GetSolutionResultsDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	results, err := h.services.SolutionResultService.ResultsBySolutionID(c.Request.Context(), dto.SolutionID)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) solutionCode(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.GetSolutionCodeDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	sol, err := h.services.SolutionService.SolutionByID(c.Request.Context(), dto.SolutionID)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, sol.Code)
}

func (h *Handler) getAvailableSolutionStatuses(c *gin.Context) {
	ss, err := h.services.SolutionManager.GetAvailableSolutionStatuses()
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, ss)
}
