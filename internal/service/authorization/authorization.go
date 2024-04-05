package authorization

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/simple_auth"
)

const (
	claimsUserKey = "user"
)

type (
	Service struct {
		config     *config.Config
		authorizer *simple_auth.Authorizer
		repository AuthorizationRepo
	}
)

func NewService(conf *config.Config, repository AuthorizationRepo) *Service {
	authorizer := simple_auth.NewAuthorizer(
		conf.Auth.AccessTokenExpTime,
		conf.Auth.RefreshTokenExpTime,
		conf.Auth.Secret,
	)

	return &Service{
		config:     conf,
		authorizer: authorizer,
		repository: repository,
	}
}

func (s *Service) Register(ctx context.Context, dto domain.CreateUserDTO) (user domain.User, err error) {
	passHash, err := simple_auth.HashPassword(dto.Password)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "Register auth service")
	}

	entity := domain.CreateUserEntity{
		Username:     dto.Username,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		PasswordHash: passHash,
	}

	user, err = s.repository.CreateUser(ctx, entity)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "Register auth service")
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, dto domain.LoginDTO) (tokens simple_auth.Tokens, err error) {
	user, err := s.repository.UserByUsername(ctx, dto.Username)
	if err != nil {
		return simple_auth.Tokens{}, errors.Wrap(err, "Login auth service")
	}

	if !simple_auth.CheckPasswordHash(dto.Password, user.PasswordHash) {
		return simple_auth.Tokens{}, errors.Wrap(errors.New("Invalid password"), "Login auth service")
	}

	tokens, err = s.createTokens(user)
	if err != nil {
		return simple_auth.Tokens{}, errors.Wrap(err, "Login auth service")
	}

	return tokens, nil
}

func (s *Service) createTokens(user domain.User) (tokens simple_auth.Tokens, err error) {
	tokens, err = s.authorizer.CreateAuthTokens(map[string]interface{}{
		claimsUserKey: user,
	})
	if err != nil {
		return simple_auth.Tokens{}, errors.Wrap(err, "createTokens auth service")
	}

	return tokens, nil
}

func (s *Service) ParseUserFromToken(ctx context.Context, accessToken string) (user domain.User, err error) {
	claims, err := s.authorizer.ValidateToken(accessToken)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "ParseUserFromToken auth service")
	}

	userValue, ok := claims[claimsUserKey]
	if !ok {
		return domain.User{}, errors.Wrap(
			errors.New("user not found in token"),
			"ParseUserFromToken auth service",
		)
	}

	err = mapstructure.Decode(userValue, &user)
	if err != nil {
		return domain.User{}, errors.Wrap(
			errors.New("cannot parse user field in token"),
			"ParseUserFromToken auth service",
		)
	}

	return user, nil
}

func (s *Service) RefreshTokens(ctx context.Context, dto domain.RefreshTokenDTO) (tokens simple_auth.Tokens, err error) {
	userFromToken, err := s.ParseUserFromToken(ctx, dto.RefreshToken)
	if err != nil {
		return simple_auth.Tokens{}, errors.Wrap(err, "RefreshTokens auth service")
	}

	user, err := s.repository.UserByID(ctx, userFromToken.ID)
	if err != nil {
		return simple_auth.Tokens{}, errors.Wrap(err, "RefreshTokens auth service")
	}

	tokens, err = s.createTokens(user)
	if err != nil {
		return simple_auth.Tokens{}, errors.Wrap(err, "RefreshTokens auth service")
	}

	return tokens, nil

}

func (s *Service) UserByID(ctx context.Context, id string) (user domain.User, err error) {
	user, err = s.repository.UserByID(ctx, id)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "UserByID auth service")
	}

	return user, nil
}

func (s *Service) Users(ctx context.Context) ([]domain.User, error) {
	users, err := s.repository.Users(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Users auth service")
	}

	return users, nil
}

func (s *Service) ChangeUserAdminStatus(ctx context.Context, dto domain.ChangeUserAdminPermissionDTO) error {
	err := s.repository.ChangeUserAdminStatus(ctx, dto)
	if err != nil {
		return errors.Wrap(err, "ChangeUserAdminStatus auth service")
	}

	return nil
}
