package server

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"lcode/config"
	"lcode/internal/handler"
	"log"
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

	//router.Use(
	//	requestid.New(
	//		requestid.WithGenerator(func() string {
	//			return "test"
	//		}),
	//		requestid.WithCustomHeaderStrKey("CJVFX-PM"),
	//	),
	//)

	router.Use(loggerMiddleware(logger))

	pprof.Register(router)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusTeapot, gin.H{"code": "URL IS INVALID", "message": "URL IS INVALID"})
	})

	wsServer := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	router.GET("/socket.io/*any", gin.WrapH(wsServer))
	router.POST("/socket.io/*any", gin.WrapH(wsServer))

	// ws handlers
	h.WS.General.Register(wsServer)
	h.WS.Comment.Register(wsServer)
	h.WS.Version.Register(wsServer)
	h.WS.Activity.Register(wsServer)

	// http handlers
	h.HTTP.General.Register(&general.Middlewares{Access: middlewares.Access, General: middlewares.General}, router)
	h.HTTP.Version.Register(&version.Middlewares{Access: middlewares.Access, Version: middlewares.Version}, router)
	h.HTTP.Comment.Register(&comment.Middlewares{Access: middlewares.Access, Comment: middlewares.Comment}, router)
	h.HTTP.Pipeline.Register(&pipeline.Middlewares{Access: middlewares.Access, Pipeline: middlewares.Pipeline}, router)
	h.HTTP.Activity.Register(middlewares.Access, router)

	go func() {
		if err := wsServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

	return &Server{
		config:    config,
		GinRouter: router,
	}
}
