package token

import (
	"errors"
	"net"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrRefreshTokenUsed    = errors.New("refresh token used")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserIDRequired      = errors.New("user id required")
	ErrPairIDRequired      = errors.New("pair id required")
	ErrIPRequired          = errors.New("ip required")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrBadRequest          = errors.New("bad request")
	ErrInvalidToken        = errors.New("invalid token")
)

type PairParams struct {
	IP                   net.IP
	UserID               uuid.UUID // GUID = UUID
	AccessTokenOptClaims jwt.MapClaims
}

type RefreshTokensParams struct {
	IP                   net.IP
	AccessTokenOptClaims jwt.MapClaims
	UserID               uuid.UUID
	RefreshToken         string
	AccessToken          string
}
