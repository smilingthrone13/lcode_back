package auth

import (
	"context"
	"lcode/internal/domain"
	"lcode/pkg/simple_auth"
)

type (
	Authorization interface {
		ParseUserFromToken(ctx context.Context, accessToken string) (user domain.User, err error)

		Register(ctx context.Context, dto domain.CreateUserDTO) (user domain.User, err error)
		Login(ctx context.Context, dto domain.LoginDTO) (tokens simple_auth.Tokens, err error)
		RefreshTokens(ctx context.Context, dto domain.RefreshTokenDTO) (tokens simple_auth.Tokens, err error)

		UserByID(ctx context.Context, id string) (user domain.User, err error)
		Users(ctx context.Context) ([]domain.User, error)
		UpdateUser(ctx context.Context, dto domain.UpdateUserDTO) (user domain.User, err error)
	}

	AuthorizationRepo interface {
		CreateUser(ctx context.Context, dto domain.CreateUserEntity) (user domain.User, err error)
		UpdateUser(ctx context.Context, dto domain.UpdateUserEntity) (user domain.User, err error)

		UserByUsername(ctx context.Context, username string) (user domain.User, err error)
		UserByID(ctx context.Context, id string) (user domain.User, err error)
		Users(ctx context.Context) ([]domain.User, error)
	}
)
