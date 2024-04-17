package app

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/internal/handler"
	"lcode/internal/handler/middleware"
	"lcode/internal/infra/database"
	"lcode/internal/infra/repository"
	"lcode/internal/infra/webapi"
	"lcode/internal/manager"
	"lcode/internal/server"
	"lcode/internal/service"
	"lcode/pkg/logger"
	"lcode/pkg/postgres"
	"lcode/pkg/struct_errors"
	"log"
	"log/slog"
	"os"
	"time"
)

const (
	logBufferSize    = 1024 * 200
	logBufferTimeout = time.Second * 3
	logFilePath      = "./lcode.log"
)

type App struct {
	Server *server.Server
	l      *slog.Logger
	cfg    *config.Config
}

func Init(cfg *config.Config) *App {
	_, l, err := logger.New(&logger.Options{
		LogFilePath:        logFilePath,
		BufferSize:         logBufferSize,
		BufferFlushTimeout: logBufferTimeout,
		DebugMode:          cfg.IsDebug,
	})

	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(cfg.Files.MainFolder, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}

	transactionProvider := postgres.NewTransactionProvider(db.GetDb())

	// init infrastructure
	repos := repository.New(&repository.InitParams{Config: cfg, DB: db})

	apis := webapi.New(&webapi.InitParams{Config: cfg})

	services := service.New(
		&service.InitParams{
			Config:             cfg,
			Logger:             l,
			TransactionManager: transactionProvider,
		},
		repos,
	)

	managers := manager.New(
		&manager.InitParams{
			Config:             cfg,
			Logger:             l,
			TransactionManager: transactionProvider,
		},
		services,
		apis,
	)

	err = setDefaultData(services, managers)
	if err != nil {
		log.Fatal(err)
	}

	handlers := handler.New(
		&handler.InitParams{
			Config:             cfg,
			Logger:             l,
			TransactionManager: transactionProvider,
		},
		services,
		managers,
	)

	middlewares := middleware.New(&middleware.InitParams{Config: cfg, Logger: l}, services, managers)

	s := server.NewServer(cfg, l, handlers, middlewares)

	return &App{Server: s, l: l, cfg: cfg}
}

func (a *App) Run() {
	srvAddr := fmt.Sprintf("%s:%s", a.cfg.HTTP.Host, a.cfg.HTTP.Port)

	if !a.cfg.TLS.Enabled {
		a.l.Info(fmt.Sprintf("server starting on http://%s", srvAddr))

		if err := a.Server.GinRouter.Run(srvAddr); err != nil {
			a.l.Error("failed run app: ", slog.String("err", err.Error()))
		}

		return
	}

	a.l.Info(fmt.Sprintf("server starting on https://%s", srvAddr))

	err := a.Server.GinRouter.RunTLS(srvAddr, a.cfg.TLS.CertFile, a.cfg.TLS.KeyFile)
	if err != nil {
		a.l.Error("failed run app: ", slog.String("err", err.Error()))
	}
}

func setDefaultData(s *service.Services, m *manager.Managers) error {
	// create default Admin user
	uCreateInput := domain.CreateUserDTO{
		Email:     "admin@admin",
		Username:  "admin",
		FirstName: "Admin",
		LastName:  "Admin",
		Password:  "admin",
	}

	u, err := m.UserManager.Register(context.Background(), uCreateInput)
	if err != nil {
		var errExist *struct_errors.ErrExist
		if ok := errors.As(err, &errExist); !ok {
			return errors.Wrap(err, "setDefaultData app init:")
		}

		u, err = s.Auth.UserByUsername(context.Background(), uCreateInput.Username)
		if err != nil {
			return errors.Wrap(err, "setDefaultData app init:")
		}
	}

	isAdmin := true
	uUpdateInput := domain.UpdateUserDTO{
		UserID:  u.ID,
		IsAdmin: &isAdmin,
	}

	u, err = m.UserManager.UpdateUser(context.Background(), uUpdateInput)
	if err != nil {
		return errors.Wrap(err, "setDefaultData app init:")
	}

	// create default practice article
	err = s.Article.CreateDefault(context.Background(), u)
	if err != nil {
		return errors.Wrap(err, "setDefaultData app init:")
	}

	return nil
}
