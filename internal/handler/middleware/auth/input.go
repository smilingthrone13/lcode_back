package auth

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/filesystem"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
)

type (
	Services struct {
	}

	Middleware struct {
		cfg        *config.Config
		logger     *slog.Logger
		services   *Services
		filesystem *filesystem.FileSystem
	}
)

func New(cfg *config.Config, logger *slog.Logger, services *Services) *Middleware {
	return &Middleware{
		cfg:        cfg,
		logger:     logger,
		services:   services,
		filesystem: &filesystem.FileSystem{},
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
	IsAdmin *bool `json:"is_admin"`
}

func (m *Middleware) ValidateChangeUserPermissionInput(c *gin.Context) {
	dto := domain.UpdateUserDTO{
		UserID: c.Param("user_id"),
	}

	var inp changeUserAdminPermissionInput

	err := c.ShouldBindJSON(&inp)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if inp.IsAdmin == nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "is_admin is required")

		return
	}

	dto.IsAdmin = inp.IsAdmin

	c.Set(domain.DtoCtxKey, dto)
}

type changeUserProfileInput struct {
	Email     *string `json:"email" binding:"omitempty,email,max=100"`
	Username  *string `json:"username" binding:"omitempty,min=3,max=50"`
	FirstName *string `json:"first_name" binding:"omitempty,min=2,max=50"`
	LastName  *string `json:"last_name" binding:"omitempty,min=2,max=50"`
	Password  *string `json:"password" binding:"omitempty,min=5,max=50"`
}

func (i *changeUserProfileInput) IsHaveUpdates() bool {
	if i.Email != nil || i.Username != nil || i.FirstName != nil || i.LastName != nil || i.Password != nil {
		return true
	}

	return false
}

func (m *Middleware) ValidateChangeUserProfileInput(c *gin.Context) {
	dto := domain.UpdateUserDTO{
		UserID: c.Param("user_id"),
	}

	var inp changeUserProfileInput

	err := c.ShouldBindJSON(&inp)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if !inp.IsHaveUpdates() {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "no updates")

		return
	}

	dto.Email = inp.Email
	dto.Username = inp.Username
	dto.FirstName = inp.FirstName
	dto.LastName = inp.LastName
	dto.Password = inp.Password

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateUploadAvatarInput(c *gin.Context) {
	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if c.Request.ContentLength <= 0 || c.Request.ContentLength > m.cfg.Files.UserAvatarMaxSize {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, "Body size exceeds limits")

		return
	}

	fullFileName := c.Param("file_name")

	name, ext, err := m.filesystem.ParseFileName(fullFileName)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	dto := domain.UploadUserAvatarDTO{
		Media:        c.Request.Body,
		FullFileName: fullFileName,
		Name:         name,
		Extension:    ext,
		MediaType:    domain.PictureMedia,
		User:         user,
	}

	c.Set(domain.DtoCtxKey, dto)
}

func (m *Middleware) ValidateDeleteAvatarInput(c *gin.Context) {
	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	dto := domain.DeleteUserAvatarDTO{
		User: user,
	}

	c.Set(domain.DtoCtxKey, dto)
}
