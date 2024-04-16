package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	accessMiddleware "lcode/internal/handler/middleware/access"
	authMiddleware "lcode/internal/handler/middleware/auth"
	"lcode/internal/manager/user_manager"
	"lcode/internal/service/auth"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"lcode/pkg/struct_errors"
	"log/slog"
	"net/http"
)

type (
	Middlewares struct {
		Access *accessMiddleware.Middleware
		Auth   *authMiddleware.Middleware
	}

	Services struct {
		Auth        auth.Authorization
		UserManager user_manager.UserManager
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
	authGroup := httpServer.Group("/auth")
	{
		authGroup.POST("/register", middlewares.Auth.ValidateRegisterInput, h.register)

		authGroup.POST("/login", middlewares.Auth.ValidateLoginInput, h.login)

		authGroup.POST("/refresh_tokens", middlewares.Auth.ValidateRefreshTokenInput, h.refreshToken)

		authGroup.GET("/my_info", middlewares.Access.UserIdentity, h.getMyInfo)

		usersGroup := authGroup.Group(
			"/users",
			middlewares.Access.UserIdentity,
		)
		{
			usersGroup.GET("", h.users)

			usersGroup.POST(
				"/upload_avatar/:file_name",
				middlewares.Auth.ValidateUploadAvatarInput,
				h.uploadAvatar,
			)

			usersGroup.DELETE(
				"/delete_avatar",
				middlewares.Auth.ValidateDeleteAvatarInput,
				h.deleteAvatar,
			)

			usersGroup.PATCH(
				"/:user_id/change_admin_permission",
				middlewares.Auth.CheckAdminAccess,
				middlewares.Auth.ValidateChangeUserPermissionInput,
				h.changeAdminPermission,
			)

			usersGroup.PATCH(
				"/:user_id/update_profile",
				middlewares.Auth.ValidateChangeUserProfileInput,
				middlewares.Auth.CheckUpdateProfilePermission,
				h.updateProfile,
			)

			usersGroup.GET(
				"/:user_id/avatar",
				h.avatarFile,
			)

			usersGroup.GET(
				"/:user_id/avatar/thumbnail",
				h.avatarThumbnailFile,
			)
		}
	}
}

func (h *Handler) register(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.CreateUserDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	user, err := h.services.UserManager.Register(c.Request.Context(), dto)
	if err != nil {
		var errExist *struct_errors.ErrExist
		if errors.As(err, &errExist) {
			http_helper.NewErrorResponse(c, http.StatusConflict, errExist.Msg)

			return
		}

		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) login(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.LoginDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	tokens, err := h.services.UserManager.Login(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) refreshToken(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.RefreshTokenDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	tokens, err := h.services.Auth.RefreshTokens(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) getMyInfo(c *gin.Context) {
	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) users(c *gin.Context) {
	users, err := h.services.UserManager.Users(c.Request.Context())
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) changeAdminPermission(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.UpdateUserDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	user, err := h.services.UserManager.UpdateUser(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) updateProfile(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.UpdateUserDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	user, err := h.services.UserManager.UpdateUser(c.Request.Context(), dto)
	if err != nil {
		var errExist *struct_errors.ErrExist
		if errors.As(err, &errExist) {
			http_helper.NewErrorResponse(c, http.StatusConflict, errExist.Msg)

			return
		}

		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) uploadAvatar(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.UploadUserAvatarDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	thumbnailPath, err := h.services.UserManager.UploadUserAvatar(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.File(thumbnailPath)
}

func (h *Handler) deleteAvatar(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.DeleteUserAvatarDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.UserManager.DeleteUserAvatar(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, "ok")
}

func (h *Handler) avatarFile(c *gin.Context) {
	p, err := h.services.UserManager.AvatarPath(c.Request.Context(), c.Param("user_id"))
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.Header("Cache-Control", "no-cache")

	c.File(p)
}

func (h *Handler) avatarThumbnailFile(c *gin.Context) {
	p, err := h.services.UserManager.AvatarThumbnailPath(c.Request.Context(), c.Param("user_id"))
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.Header("Cache-Control", "no-cache")

	c.File(p)
}
