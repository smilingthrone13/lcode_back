package solution

import (
	"github.com/gin-gonic/gin"
	"lcode/internal/domain"
	"lcode/pkg/gin_helpers"
	"lcode/pkg/http_lib/http_helper"
	"lcode/pkg/struct_errors"
	"net/http"
)

func (m *Middleware) CheckSolutionAccess(c *gin.Context) {
	dto, err := gin_helpers.GetValueFromGinCtx[domain.IGetSolutionDTO](c, domain.DtoCtxKey)
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	sol, err := m.services.Solution.SolutionByID(c.Request.Context(), dto.GetSolutionID())
	if err != nil {
		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if sol.UserID != dto.GetUser().ID {
		err = struct_errors.NewBaseErr("No access rights", nil)

		http_helper.NewErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}
}
