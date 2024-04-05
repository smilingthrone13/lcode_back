package access

import (
	"context"
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/http_lib/http_helper"
	"log/slog"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

type (
	AuthService interface {
		ParseUserFromToken(ctx context.Context, accessToken string) (user domain.User, err error)
	}

	Services struct {
		Auth AuthService
	}
)

type Middleware struct {
	cfg      *config.Config
	logger   *slog.Logger
	services *Services
}

func New(cfg *config.Config, logger *slog.Logger, services *Services) *Middleware {
	return &Middleware{
		cfg:      cfg,
		logger:   logger,
		services: services,
	}
}

func (m *Middleware) UserIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		http_helper.NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")

		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		http_helper.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")

		return
	}
	if len(headerParts[1]) == 0 {
		http_helper.NewErrorResponse(c, http.StatusUnauthorized, "token is empty")

		return
	}

	user, err := m.services.Auth.ParseUserFromToken(c.Request.Context(), headerParts[1])
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusUnauthorized, err.Error())

		return
	}

	c.Set(domain.UserCtxKey, user)
}
