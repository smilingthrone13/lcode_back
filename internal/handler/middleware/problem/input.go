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

	if dto.Input.Task.MemoryLimit == 0 {
		dto.Input.Task.MemoryLimit = m.cfg.JudgeConfig.DefaultMemoryLimitKB
	}

	if dto.Input.Task.RuntimeLimit == 0 {
		dto.Input.Task.RuntimeLimit = m.cfg.JudgeConfig.DefaultTimeLimitSec
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

	if dto.Input.Name == nil && dto.Input.Category == nil &&
		dto.Input.Description == nil && dto.Input.Difficulty == nil &&
		dto.Input.MemoryLimit == nil && dto.Input.RuntimeLimit == nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "No update data provided")

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

	if dto.Input.Template == "" || dto.Input.Wrapper == "" || dto.Input.LanguageID == 0 {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Invalid input")

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

	if dto.Input.Wrapper == nil && dto.Input.Template == nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "No update data provided")

		return
	}

	dto.TemplateID = c.Param("template_id")
	dto.TaskID = c.Param("task_id")
	if dto.TemplateID == "" || dto.TaskID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Task ID and Template ID are required")

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

	if dto.Input.Input == "" || dto.Input.Output == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Invalid input")

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

	if dto.Input.Input == nil && dto.Input.Output == nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "No update data provided")

		return
	}

	dto.CaseID = c.Param("case_id")
	dto.TaskID = c.Param("task_id")
	if dto.CaseID == "" || dto.TaskID == "" {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Task ID and TestCase ID are required")

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

	if err := c.ShouldBindJSON(&inp); err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	categories, ok := c.GetQueryArray("category")
	if !ok {
		categories = []string{}
	}

	difficulties, ok := c.GetQueryArray("difficulty")
	if !ok {
		difficulties = []string{}
	}

	filter := domain.TaskFilter{
		Search:       c.Query("search"),
		Categories:   categories,
		Difficulties: difficulties,
	}

	data := domain.TaskParams{
		Filter:     filter,
		Sort:       inp.Sort,
		Pagination: inp.Pagination,
	}

	c.Set(domain.DtoCtxKey, domain.TaskParamsDTO{Input: data})
}
