package auth

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
)

type (
	Services struct {
	}

	Middleware struct {
		cfg      *config.Config
		logger   *slog.Logger
		services *Services
	}
)

func New(cfg *config.Config, logger *slog.Logger, services *Services) *Middleware {
	return &Middleware{
		cfg:      cfg,
		logger:   logger,
		services: services,
	}
}

func (m *Middleware) ValidateRegisterInput(c *gin.Context) {
	var dto domain.CreateUserDTO

	err := c.ShouldBindJSON(&dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateLoginInput(c *gin.Context) {
	var dto domain.LoginDTO

	err := c.ShouldBindJSON(&dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateRefreshTokenInput(c *gin.Context) {
	var dto domain.RefreshTokenDTO

	err := c.ShouldBindJSON(&dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.Set(domain.DtoCtxKey, dto)
}

type changeUserAdminPermissionInput struct {
	IsAdmin bool `json:"is_admin"`
}

func (m *Middleware) ValidateChangeUserPermissionInput(c *gin.Context) {
	dto := domain.ChangeUserAdminPermissionDTO{
		UserID: c.Param("user_id"),
	}

	var inp changeUserAdminPermissionInput

	err := c.ShouldBindJSON(&inp)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	dto.IsAdmin = inp.IsAdmin

	c.Set(domain.DtoCtxKey, dto)
}
