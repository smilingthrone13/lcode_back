package auth

import (
	"github.com/gin-gonic/gin"
	"lcode/internal/domain"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"net/http"
)

func (m *Middleware) CheckAdminAccess(c *gin.Context) {
	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if !user.IsAdmin {
		http_helper.NewErrorResponse(c, http.StatusForbidden, "user is not admin")

		return
	}
}

func (m *Middleware) CheckUpdateProfilePermission(c *gin.Context) {
	user, err := gin_helpers.GetValueFromGinCtx[domain.User](c, domain.UserCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	dto, err := gin_helpers.GetValueFromGinCtx[domain.UpdateUserDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if dto.UserID != user.ID && !user.IsAdmin {
		http_helper.NewErrorResponse(c, http.StatusForbidden, "no permissions")

		return
	}

}
