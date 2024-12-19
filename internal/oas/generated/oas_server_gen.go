// Code generated by ogen, DO NOT EDIT.

package gen

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// CreateTokens implements CreateTokens operation.
	//
	// Получение новой пары токенов.
	//
	// POST /users/{user_id}/tokens
	CreateTokens(ctx context.Context, req OptCreateTokensRequest, params CreateTokensParams) (CreateTokensRes, error)
	// RefreshTokens implements RefreshTokens operation.
	//
	// Получение новой пары токенов по рефреш токену.
	//
	// PUT /users/{user_id}/tokens
	RefreshTokens(ctx context.Context, req *RefreshTokensReq, params RefreshTokensParams) (RefreshTokensRes, error)
	// NewError creates *ErrorResponseStatusCode from error returned by handler.
	//
	// Used for common default response.
	NewError(ctx context.Context, err error) *ErrorResponseStatusCode
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h Handler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		baseServer: s,
	}, nil
}