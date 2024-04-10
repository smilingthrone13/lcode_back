package problem

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	accessMiddleware "lcode/internal/handler/middleware/access"
	authMiddleware "lcode/internal/handler/middleware/auth"
	problemMiddleware "lcode/internal/handler/middleware/problem"
	problemManager "lcode/internal/manager/problem_manager"
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
		Problem *problemMiddleware.Middleware
	}

	Managers struct {
		Problem problemManager.ProblemManager
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
	problemGroup := httpServer.Group("/problems", middlewares.Access.UserIdentity)
	{
		problemGroup.GET(
			"/:task_id",
			middlewares.Problem.ValidateFullProblemByTaskIDInput,
			h.getProblem,
		)
		problemGroup.GET(
			"/list",
			middlewares.Problem.ValidateTaskListByParamsInput,
			h.getTasksList,
		)
		problemGroup.GET(
			"/available_attributes",
			h.getAvailableTaskAttributes,
		)

		taskGroup := problemGroup.Group("", middlewares.Auth.CheckAdminAccess)
		{
			taskGroup.POST(
				"/",
				middlewares.Problem.ValidateCreateProblemInput,
				h.createProblem,
			)
			taskGroup.PATCH(
				"/:task_id",
				middlewares.Problem.ValidateUpdateProblemTaskInput,
				h.updateProblemTask,
			)
			taskGroup.DELETE(
				"/:task_id",
				middlewares.Problem.ValidateDeleteProblemInput,
				h.deleteProblem,
			)
		}

		templateGroup := problemGroup.Group("/:task_id/template", middlewares.Auth.CheckAdminAccess)
		{
			templateGroup.POST(
				"/",
				middlewares.Problem.ValidateCreateProblemTaskTemplateInput,
				h.createProblemTaskTemplate,
			)
			templateGroup.PATCH(
				"/:template_id",
				middlewares.Problem.ValidateUpdateProblemTaskTemplateInput,
				h.updateProblemTaskTemplate,
			)
			templateGroup.DELETE(
				"/:template_id",
				middlewares.Problem.ValidateDeleteProblemTaskTemplateInput,
				h.deleteProblemTaskTemplate,
			)
		}

		testCaseGroup := problemGroup.Group("/:task_id/testcase", middlewares.Auth.CheckAdminAccess)
		{
			testCaseGroup.POST(
				"/",
				middlewares.Problem.ValidateCreateProblemTestCaseInput,
				h.createProblemTestCase,
			)
			testCaseGroup.PATCH(
				"/:case_id",
				middlewares.Problem.ValidateUpdateProblemTestCaseInput,
				h.updateProblemTestCase,
			)
			testCaseGroup.DELETE(
				"/:case_id",
				middlewares.Problem.ValidateDeleteProblemTestCaseInput,
				h.deleteProblemTestCase,
			)
		}
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
		var errExist *struct_errors.ErrExist
		if errors.As(err, &errExist) {
			http_helper.NewErrorResponse(c, http.StatusConflict, errExist.Msg)

			return
		}

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

	c.JSON(http.StatusOK, map[string]string{"message": "Successful operation"})
}

func (h *Handler) createProblemTaskTemplate(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.TaskTemplateCreateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	problem, err := h.managers.Problem.CreateProblemTaskTemplate(c.Request.Context(), dto.TaskID, dto.Input)
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

	c.JSON(http.StatusOK, map[string]string{"message": "Successful operation"})
}

func (h *Handler) createProblemTestCase(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.TestCaseCreateDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	problem, err := h.managers.Problem.CreateProblemTestCase(c.Request.Context(), dto.TaskID, dto.Input)
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

	c.JSON(http.StatusOK, map[string]string{"message": "Successful operation"})
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

func (h *Handler) getAvailableTaskAttributes(c *gin.Context) {
	ta, err := h.managers.Problem.GetAvailableTaskAttributes(c.Request.Context())
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, ta)
}
