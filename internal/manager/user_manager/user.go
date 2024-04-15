package user_manager

import (
	"context"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/internal/service/auth"
	"lcode/internal/service/user_fs"
	"lcode/pkg/filesystem"
	"lcode/pkg/postgres"
	"lcode/pkg/simple_auth"
	"log"
	"log/slog"
	"os"
	"path"
)

type (
	Services struct {
		Auth   auth.Authorization
		UserFS user_fs.UserFS
	}

	Manager struct {
		cfg                *config.Config
		logger             *slog.Logger
		transactionManager *postgres.TransactionProvider
		services           *Services
		fileSystem         *filesystem.FileSystem
	}
)

func New(
	cfg *config.Config,
	logger *slog.Logger,
	transactionManager *postgres.TransactionProvider,
	services *Services,
) *Manager {

	err := os.MkdirAll(path.Join(cfg.Files.MainFolder, domain.UsersFolder), os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}

	return &Manager{
		cfg:                cfg,
		logger:             logger,
		transactionManager: transactionManager,
		services:           services,
		fileSystem:         &filesystem.FileSystem{},
	}
}

func (m *Manager) Register(
	ctx context.Context,
	dto domain.CreateUserDTO,
) (user domain.User, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "Register user manager")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	user, err = m.services.Auth.Register(ctx, dto)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "Register user manager")
	}

	err = m.services.UserFS.MakeUserDir(ctx, user)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "Register user manager")
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.User{}, errors.Wrap(err, "Register user manager")
	}

	return user, nil
}

func (m *Manager) Login(ctx context.Context, dto domain.LoginDTO) (tokens simple_auth.Tokens, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return simple_auth.Tokens{}, errors.Wrap(err, "Login user manager")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	tokens, err = m.services.Auth.Login(ctx, dto)
	if err != nil {
		return simple_auth.Tokens{}, errors.Wrap(err, "Login user manager")
	}

	if err = tx.Commit(ctx); err != nil {
		return simple_auth.Tokens{}, errors.Wrap(err, "Login user manager")
	}

	return tokens, nil
}

func (m *Manager) UserByID(ctx context.Context, id string) (user domain.User, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "UserByID user manager")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	user, err = m.services.Auth.UserByID(ctx, id)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "UserByID user manager")
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.User{}, errors.Wrap(err, "UserByID user manager")
	}

	return user, nil
}

func (m *Manager) Users(ctx context.Context) ([]domain.User, error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Users user manager")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	users, err := m.services.Auth.Users(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Users user manager")
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "Users user manager")
	}

	return users, nil
}

func (m *Manager) UpdateUser(ctx context.Context, dto domain.UpdateUserDTO) (user domain.User, err error) {
	tx, err := m.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "UpdateUser user manager")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	user, err = m.services.Auth.UpdateUser(ctx, dto)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "UpdateUser user manager")
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.User{}, errors.Wrap(err, "UpdateUser user manager")
	}

	return user, nil
}

func (m *Manager) UploadUserAvatar(
	ctx context.Context,
	dto domain.UploadUserAvatarDTO,
) (thumbnailPath string, err error) {
	err = m.services.UserFS.DeleteAvatar(ctx, dto.User.ID)
	if err != nil {
		return "", errors.Wrap(err, "UploadUserAvatar user manager")
	}

	_, thumbnailPath, err = m.services.UserFS.CreateAvatar(ctx, dto)
	if err != nil {
		return "", errors.Wrap(err, "UploadUserAvatar user manager")
	}

	return thumbnailPath, nil
}

func (m *Manager) DeleteUserAvatar(ctx context.Context, dto domain.DeleteUserAvatarDTO) error {
	err := m.services.UserFS.DeleteAvatar(ctx, dto.User.ID)
	if err != nil {
		return errors.Wrap(err, "DeleteUserAvatar user manager")
	}

	return nil
}

func (m *Manager) AvatarPath(ctx context.Context, userID string) (p string, err error) {
	origPath, err := m.services.UserFS.AvatarPath(ctx, userID)
	if err != nil {
		return "", errors.Wrap(err, "AvatarPath user manager")
	}

	return origPath, nil
}

func (m *Manager) AvatarThumbnailPath(ctx context.Context, userID string) (p string, err error) {
	thumbnailPath, err := m.services.UserFS.AvatarThumbnailPath(ctx, userID)
	if err != nil {
		return "", errors.Wrap(err, "AvatarThumbnailPath user manager")
	}

	return thumbnailPath, nil
}
