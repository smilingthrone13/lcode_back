package problem

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/internal/manager/problem_manager"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
)

type (
	Managers struct {
		Problem *problem_manager.Manager
	}

	Middleware struct {
		cfg      *config.Config
		logger   *slog.Logger
		managers *Managers
	}
)

func New(cfg *config.Config, logger *slog.Logger, managers *Managers) *Middleware {
	return &Middleware{
		cfg:      cfg,
		logger:   logger,
		managers: managers,
	}
}

func (m *Middleware) ValidateCreateProblemInput(c *gin.Context) {
	var dto domain.ProblemCreateDTO

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if dto.Input.Task.Name == "" ||
		dto.Input.Task.Category == "" ||
		dto.Input.Task.Difficulty == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Invalid input")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateUpdateProblemTaskInput(c *gin.Context) {
	var dto domain.TaskUpdateDTO

	err := c.ShouldBindJSON(&dto.Input)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	dto.TaskID = c.Param("task_id")
	if dto.TaskID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Task ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateDeleteProblemInput(c *gin.Context) {
	dto := domain.ProblemDeleteDTO{
		TaskID: c.Param("task_id"),
	}

	if dto.TaskID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Task ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateCreateProblemTaskTemplateInput(c *gin.Context) {
	var dto domain.TaskTemplateCreateDTO

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	dto.TaskID = c.Param("task_id")
	if dto.TaskID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Task ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateUpdateProblemTaskTemplateInput(c *gin.Context) {
	var dto domain.TaskTemplateUpdateDTO

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	dto.TemplateID = c.Param("template_id")
	if dto.TemplateID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Template ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateDeleteProblemTaskTemplateInput(c *gin.Context) {
	dto := domain.TaskTemplateDeleteDTO{
		TemplateID: c.Param("template_id"),
	}

	if dto.TemplateID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Template ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateCreateProblemTestCaseInput(c *gin.Context) {
	var dto domain.TestCaseCreateDTO

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	dto.TaskID = c.Param("task_id")
	if dto.TaskID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Task ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateUpdateProblemTestCaseInput(c *gin.Context) {
	var dto domain.TestCaseUpdateDTO

	if err := c.ShouldBindJSON(&dto.Input); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	dto.CaseID = c.Param("case_id")
	if dto.CaseID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Test Case ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateDeleteProblemTestCaseInput(c *gin.Context) {
	dto := domain.TestCaseDeleteDTO{
		CaseID: c.Param("case_id"),
	}

	if dto.CaseID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Test Case ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateFullProblemByTaskIDInput(c *gin.Context) {
	dto := domain.GetProblemDTO{
		TaskID: c.Param("task_id"),
	}

	if dto.TaskID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Task ID is required")

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateTaskListByParamsInput(c *gin.Context) {
	var inp domain.TaskParamsInput

	if err := c.ShouldBindQuery(&inp); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	filter := domain.TaskFilter{
		Search:       c.Query("search"),
		Categories:   c.QueryArray("categories"),
		Difficulties: c.QueryArray("difficulties"),
	}

	data := domain.TaskParams{
		Filter:     filter,
		Sort:       inp.Sort,
		Pagination: inp.Pagination,
	}

	c.Set(domain.DtoCtxKey, domain.TaskParamsDTO{Input: data})
}
