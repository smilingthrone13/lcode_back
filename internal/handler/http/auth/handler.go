package auth

import (
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	accessMiddleware "lcode/internal/handler/middleware/access"
	authMiddleware "lcode/internal/handler/middleware/auth"
	"lcode/internal/service/authorization"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
)

type (
	Middlewares struct {
		Access *accessMiddleware.Middleware
		Auth   *authMiddleware.Middleware
	}

	Services struct {
		Auth authorization.Authorization
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

		// todo: add middleware for admin
		authGroup.GET("/users", middlewares.Access.UserIdentity, middlewares.Auth.CheckAdminAccess, h.users)

		// todo: add middleware for admin
		usersGroup := authGroup.Group("/user", middlewares.Access.UserIdentity, middlewares.Auth.CheckAdminAccess)
		{

			usersGroup.PATCH(
				"/:user_id/change_admin_permission",
				middlewares.Auth.ValidateChangeUserPermissionInput,
				h.changeAdminPermission,
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

	user, err := h.services.Auth.Register(c.Request.Context(), dto)
	if err != nil {
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

	tokens, err := h.services.Auth.Login(c.Request.Context(), dto)
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

func (h *Handler) users(c *gin.Context) {
	users, err := h.services.Auth.Users(c.Request.Context())
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) changeAdminPermission(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.ChangeUserAdminPermissionDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.Auth.ChangeUserAdminStatus(c.Request.Context(), dto)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	c.JSON(http.StatusOK, "ok")
}
