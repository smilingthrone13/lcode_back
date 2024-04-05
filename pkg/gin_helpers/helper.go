package gin_helpers

import (
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log/slog"
)

func GetValueFromGinCtx[T any](c *gin.Context, key string) (value T, err error) {
	val, ok := c.Get(key)
	if !ok {
		return value, errors.New(fmt.Sprintf("cannot find key `%s` in gin ctx", key))
	}

	value, ok = val.(T)
	if !ok {
		return value, errors.New(fmt.Sprintf("cannot type assert key `%s` in gin ctx", key))
	}

	return
}

func GetRequestLogAttr(c *gin.Context) slog.Attr {
	return slog.Group("request_params",
		slog.String("request_id", requestid.Get(c)),
		slog.String("addr", c.Request.RemoteAddr),
		slog.String("path", c.Request.URL.RequestURI()),
	)
}
