package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/handler"
	"lcode/internal/handler/http/article"
	"lcode/internal/handler/http/auth"
	"lcode/internal/handler/http/comment"
	"lcode/internal/handler/http/problem"
	"lcode/internal/handler/http/solution"
	userProgress "lcode/internal/handler/http/user_progress"
	"lcode/internal/handler/middleware"
	"log/slog"
	"net/http"
	"time"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

type Server struct {
	config    *config.Config
	GinRouter *gin.Engine
}

func loggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		logger.DebugContext(c.Request.Context(), "head of request", slog.Group(
			"params",
			slog.String("request_id", requestid.Get(c)),
			slog.String("addr", c.Request.RemoteAddr),
			slog.String("path", c.Request.URL.RequestURI()),
		))

		c.Next()

		// after request
		latency := time.Since(t)

		logger.DebugContext(c.Request.Context(), "end of request", slog.Group(
			"params",
			slog.String("request_id", requestid.Get(c)),
			slog.Int64("request_time_ms", latency.Milliseconds()),
		))
	}
}

func NewServer(
	config *config.Config,
	logger *slog.Logger,
	h *handler.Handlers,
	middlewares *middleware.Middlewares,
) *Server {
	if config.IsDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = append(cfg.AllowOrigins, config.CorsOrigins...)
	cfg.AllowCredentials = true
	cfg.AllowHeaders = append(cfg.AllowHeaders,
		"Access-Control-Allow-Headers",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers",
		"Accept",
		"X-Requested-With",
		"Authorization")
	router.Use(cors.New(cfg))

	router.Use(requestid.New())

	router.Use(loggerMiddleware(logger))

	pprof.Register(router)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusTeapot, gin.H{"code": "URL IS INVALID", "message": "URL IS INVALID"})
	})

	// http handlers
	h.HTTP.Auth.Register(
		&auth.Middlewares{
			Access: middlewares.Access,
			Auth:   middlewares.Auth,
		},
		router,
	)

	h.HTTP.Problem.Register(
		&problem.Middlewares{
			Problem: middlewares.Problem,
			Access:  middlewares.Access,
			Auth:    middlewares.Auth,
		},
		router,
	)

	h.HTTP.UserProgress.Register(
		&userProgress.Middlewares{
			Access:       middlewares.Access,
			UserProgress: middlewares.UserProgress,
		},
		router,
	)

	h.HTTP.Article.Register(
		&article.Middlewares{
			Access:  middlewares.Access,
			Auth:    middlewares.Auth,
			Article: middlewares.Article,
		},
		router,
	)

	h.HTTP.Solution.Register(
		&solution.Middlewares{
			Access:   middlewares.Access,
			Solution: middlewares.Solution,
		},
		router,
	)

	h.HTTP.Comment.Register(
		&comment.Middlewares{
			Access:  middlewares.Access,
			Comment: middlewares.Comment,
		},
		router,
	)

	return &Server{
		config:    config,
		GinRouter: router,
	}
}
