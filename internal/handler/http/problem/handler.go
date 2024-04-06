package problem

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	accessMiddleware "lcode/internal/handler/middleware/access"
	authMiddleware "lcode/internal/handler/middleware/auth"
	problemMiddleware "lcode/internal/handler/middleware/problem"
	"lcode/internal/manager/problem_manager"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
)

type (
	Middlewares struct {
		Auth    *authMiddleware.Middleware
		Access  *accessMiddleware.Middleware
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
	problemGroup := httpServer.Group("/problem", middlewares.Access.UserIdentity)
	{
		problemGroup.GET(
			"/",
			middlewares.Problem.ValidateFullProblemByTaskIDInput,
			h.getProblem,
		)
		problemGroup.GET(
			"/list",
			middlewares.Problem.ValidateTaskListByParamsInput,
			h.getTasksList,
		)

		problemGroup.POST(
			"/",
			middlewares.Auth.CheckAdminAccess,
			middlewares.Problem.ValidateCreateProblemInput,
			h.createProblem,
		)
		problemGroup.PATCH(
			"/",
			middlewares.Auth.CheckAdminAccess,
			middlewares.Problem.ValidateUpdateProblemTaskInput,
			h.updateProblemTask,
		)
		problemGroup.DELETE(
			"/",
			middlewares.Auth.CheckAdminAccess,
			middlewares.Problem.ValidateDeleteProblemInput,
			h.deleteProblem,
		)

		problemGroup.POST(
			"/template",
			middlewares.Auth.CheckAdminAccess,
			middlewares.Problem.ValidateCreateProblemTaskTemplateInput,
			h.createProblemTaskTemplate,
		)
		problemGroup.PATCH(
			"/template",
			middlewares.Auth.CheckAdminAccess,
			middlewares.Problem.ValidateUpdateProblemTaskTemplateInput,
			h.updateProblemTaskTemplate,
		)
		problemGroup.DELETE(
			"/template",
			middlewares.Auth.CheckAdminAccess,
			middlewares.Problem.ValidateDeleteProblemTaskTemplateInput,
			h.deleteProblemTaskTemplate,
		)

		problemGroup.POST(
			"/testcase",
			middlewares.Auth.CheckAdminAccess,
			middlewares.Problem.ValidateCreateProblemTestCaseInput,
			h.createProblemTestCase,
		)
		problemGroup.PATCH(
			"/testcase",
			middlewares.Auth.CheckAdminAccess,
			middlewares.Problem.ValidateUpdateProblemTestCaseInput,
			h.updateProblemTestCase,
		)
		problemGroup.DELETE(
			"/testcase",
			middlewares.Auth.CheckAdminAccess,
			middlewares.Problem.ValidateDeleteProblemTestCaseInput,
			h.deleteProblemTestCase,
		)
	}
}

func (h *Handler) createProblem(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.ProblemCreateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	problem, err := h.managers.Problem.CreateProblem(c.Request.Context(), dto.Input)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusCreated, problem)

}

func (h *Handler) updateProblemTask(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.TaskUpdateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	problem, err := h.managers.Problem.UpdateProblemTask(c.Request.Context(), dto.TaskID, dto.Input)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, problem)
}

func (h *Handler) deleteProblem(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.ProblemDeleteDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.managers.Problem.DeleteProblem(c.Request.Context(), dto.TaskID)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) createProblemTaskTemplate(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.TaskTemplateCreateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	problem, err := h.managers.Problem.CreateProblemTaskTemplate(c.Request.Context(), dto.Input)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusCreated, problem)
}

func (h *Handler) updateProblemTaskTemplate(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.TaskTemplateUpdateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	problem, err := h.managers.Problem.UpdateProblemTaskTemplate(c.Request.Context(), dto.TemplateID, dto.Input)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, problem)
}

func (h *Handler) deleteProblemTaskTemplate(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.TaskTemplateDeleteDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.managers.Problem.DeleteProblemTaskTemplate(c.Request.Context(), dto.TemplateID)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) createProblemTestCase(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.TestCaseCreateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	problem, err := h.managers.Problem.CreateProblemTestCase(c.Request.Context(), dto.Input)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusCreated, problem)
}

func (h *Handler) updateProblemTestCase(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.TestCaseUpdateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	problem, err := h.managers.Problem.UpdateProblemTestCase(c.Request.Context(), dto.CaseID, dto.Input)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, problem)
}

func (h *Handler) deleteProblemTestCase(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.TestCaseDeleteDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.managers.Problem.DeleteProblemTestCase(c.Request.Context(), dto.CaseID)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) getProblem(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.GetProblemDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	problem, err := h.managers.Problem.FullProblemByTaskID(c.Request.Context(), dto.TaskID)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, problem)
}

func (h *Handler) getTasksList(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.TaskParamsDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	tasks, err := h.managers.Problem.TaskListByParams(c.Request.Context(), dto.Input)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, tasks)
}
