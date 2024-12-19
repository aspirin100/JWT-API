package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/aspirin100/JWT-API/internal/logger"
	"github.com/aspirin100/JWT-API/internal/middleware"
	gen "github.com/aspirin100/JWT-API/internal/oas/generated"
	"github.com/aspirin100/JWT-API/internal/token"
)

var _ gen.Handler = (*Handler)(nil)

//go:generate go run go.uber.org/mock/mockgen@latest -source=handler.go -destination=handler_mocks_test.go -package=api_test

type TokenService interface {
	CreateNewTokensPair(ctx context.Context, params *token.PairParams) (*string, *string, error)
	RefreshTokenPair(ctx context.Context, params *token.RefreshTokensParams) (*string, *string, error)
}

type Handler struct {
	TokenService TokenService
}

func (h *Handler) CreateTokens(
	ctx context.Context,
	req gen.OptCreateTokensRequest,
	params gen.CreateTokensParams,
) (gen.CreateTokensRes, error) {
	userIP, ok := middleware.GetClientIP(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to create token pair: %w", token.ErrBadRequest)
	}

	optClaims := make(jwt.MapClaims)

	if req.Set {
		for k, v := range req.Value.GetAdditionalInfo() {
			optClaims[k] = v
		}
	}

	// создается новая пара токенов
	accessToken, refreshToken, err := h.TokenService.CreateNewTokensPair(ctx, &token.PairParams{
		IP:                   userIP,
		UserID:               params.UserID,
		AccessTokenOptClaims: optClaims,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create token pair: %w", err)
	}

	return &gen.CreateTokensResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (h *Handler) RefreshTokens(
	ctx context.Context,
	req *gen.RefreshTokensReq,
	params gen.RefreshTokensParams,
) (gen.RefreshTokensRes, error) {
	accessTokenOptClaims := make(jwt.MapClaims)

	ip, ok := middleware.GetClientIP(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to refresh token pair: %w", token.ErrBadRequest)
	}

	// рефреш операция
	accessToken, refreshToken, err := h.TokenService.RefreshTokenPair(ctx, &token.RefreshTokensParams{
		IP:                   ip,
		AccessTokenOptClaims: accessTokenOptClaims,
		RefreshToken:         req.RefreshToken,
		AccessToken:          req.AccessToken,
		UserID:               params.UserID,
	})
	if err != nil {
		logger.Default().Info("error from refresh tokens")

		return nil, fmt.Errorf("failed to refresh token pair: %w", err)
	}

	return &gen.CreateTokensResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (*Handler) NewError(_ context.Context, err error) *gen.ErrorResponseStatusCode {
	resp := &gen.ErrorResponseStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: gen.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}

	switch {
	case errors.Is(err, token.ErrRefreshTokenUsed):
		resp.StatusCode = http.StatusForbidden
		resp.Response.Code = http.StatusForbidden
		resp.Response.Message = token.ErrRefreshTokenUsed.Error()
	case errors.Is(err, token.ErrRefreshTokenExpired):
		resp.StatusCode = http.StatusForbidden
		resp.Response.Code = http.StatusForbidden
		resp.Response.Message = token.ErrRefreshTokenExpired.Error()
	case errors.Is(err, token.ErrUserNotFound):
		resp.StatusCode = http.StatusNotFound
		resp.Response.Code = http.StatusNotFound
		resp.Response.Message = token.ErrUserNotFound.Error()
	case errors.Is(err, token.ErrBadRequest):
		resp.StatusCode = http.StatusBadRequest
		resp.Response.Code = http.StatusBadRequest
		resp.Response.Message = token.ErrBadRequest.Error()
	default:
	}

	return resp
}
