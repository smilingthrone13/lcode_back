package simple_auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
)

func NewAuthorizer(
	accessTokenTime,
	refreshTokenTime time.Duration,
	secret string,
) *Authorizer {
	return &Authorizer{
		accessTokenExpTime:  accessTokenTime,
		refreshTokenExpTime: refreshTokenTime,
		secret:              []byte(secret),
	}
}

type Authorizer struct {
	accessTokenExpTime  time.Duration
	refreshTokenExpTime time.Duration
	secret              []byte
}

func (a *Authorizer) CreateAuthTokens(claims map[string]interface{}) (Tokens, error) {
	if claims == nil {
		return Tokens{}, errors.New("claims not be nil")
	}

	jwtClaims := jwt.MapClaims(claims)

	jwtClaims["exp"] = time.Now().Add(a.accessTokenExpTime).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	accessToken, err := token.SignedString(a.secret)
	if err != nil {
		return Tokens{}, err
	}

	jwtClaims["exp"] = time.Now().Add(a.refreshTokenExpTime).Unix()

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	refreshToken, err := token.SignedString(a.secret)
	if err != nil {
		return Tokens{}, err
	}

	info := Tokens{
		AccessToken:     accessToken,
		AccessTokenExp:  a.accessTokenExpTime.Milliseconds(),
		RefreshToken:    refreshToken,
		RefreshTokenExp: a.refreshTokenExpTime.Milliseconds(),
	}

	return info, nil
}

func (a *Authorizer) ValidateToken(tokenString string) (claims map[string]any, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.secret, nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

type Tokens struct {
	AccessToken     string `json:"access_token"`
	AccessTokenExp  int64  `json:"access_token_exp"`
	RefreshToken    string `json:"refresh_token"`
	RefreshTokenExp int64  `json:"refresh_token_exp"`
}
