// Code generated by ogen, DO NOT EDIT.

package client

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/uri"
)

// Invoker invokes operations described by OpenAPI v3 specification.
type Invoker interface {
	// CreateTokens invokes CreateTokens operation.
	//
	// Получение новой пары токенов.
	//
	// POST /users/{user_id}/tokens
	CreateTokens(ctx context.Context, request OptCreateTokensRequest, params CreateTokensParams) (CreateTokensRes, error)
	// RefreshTokens invokes RefreshTokens operation.
	//
	// Получение новой пары токенов по рефреш токену.
	//
	// PUT /users/{user_id}/tokens
	RefreshTokens(ctx context.Context, request *RefreshTokensReq, params RefreshTokensParams) (RefreshTokensRes, error)
}

// Client implements OAS client.
type Client struct {
	serverURL *url.URL
	baseClient
}

func trimTrailingSlashes(u *url.URL) {
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawPath = strings.TrimRight(u.RawPath, "/")
}

// NewClient initializes new Client defined by OAS.
func NewClient(serverURL string, opts ...ClientOption) (*Client, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	trimTrailingSlashes(u)

	c, err := newClientConfig(opts...).baseClient()
	if err != nil {
		return nil, err
	}
	return &Client{
		serverURL:  u,
		baseClient: c,
	}, nil
}

type serverURLKey struct{}

// WithServerURL sets context key to override server URL.
func WithServerURL(ctx context.Context, u *url.URL) context.Context {
	return context.WithValue(ctx, serverURLKey{}, u)
}

func (c *Client) requestURL(ctx context.Context) *url.URL {
	u, ok := ctx.Value(serverURLKey{}).(*url.URL)
	if !ok {
		return c.serverURL
	}
	return u
}

// CreateTokens invokes CreateTokens operation.
//
// Получение новой пары токенов.
//
// POST /users/{user_id}/tokens
func (c *Client) CreateTokens(ctx context.Context, request OptCreateTokensRequest, params CreateTokensParams) (CreateTokensRes, error) {
	res, err := c.sendCreateTokens(ctx, request, params)
	return res, err
}

func (c *Client) sendCreateTokens(ctx context.Context, request OptCreateTokensRequest, params CreateTokensParams) (res CreateTokensRes, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [3]string
	pathParts[0] = "/users/"
	{
		// Encode "user_id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "user_id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.UUIDToString(params.UserID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	pathParts[2] = "/tokens"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "POST", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeCreateTokensRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeCreateTokensResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}

// RefreshTokens invokes RefreshTokens operation.
//
// Получение новой пары токенов по рефреш токену.
//
// PUT /users/{user_id}/tokens
func (c *Client) RefreshTokens(ctx context.Context, request *RefreshTokensReq, params RefreshTokensParams) (RefreshTokensRes, error) {
	res, err := c.sendRefreshTokens(ctx, request, params)
	return res, err
}

func (c *Client) sendRefreshTokens(ctx context.Context, request *RefreshTokensReq, params RefreshTokensParams) (res RefreshTokensRes, err error) {

	u := uri.Clone(c.requestURL(ctx))
	var pathParts [3]string
	pathParts[0] = "/users/"
	{
		// Encode "user_id" parameter.
		e := uri.NewPathEncoder(uri.PathEncoderConfig{
			Param:   "user_id",
			Style:   uri.PathStyleSimple,
			Explode: false,
		})
		if err := func() error {
			return e.EncodeValue(conv.UUIDToString(params.UserID))
		}(); err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		encoded, err := e.Result()
		if err != nil {
			return res, errors.Wrap(err, "encode path")
		}
		pathParts[1] = encoded
	}
	pathParts[2] = "/tokens"
	uri.AddPathParts(u, pathParts[:]...)

	r, err := ht.NewRequest(ctx, "PUT", u)
	if err != nil {
		return res, errors.Wrap(err, "create request")
	}
	if err := encodeRefreshTokensRequest(request, r); err != nil {
		return res, errors.Wrap(err, "encode request")
	}

	resp, err := c.cfg.Client.Do(r)
	if err != nil {
		return res, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	result, err := decodeRefreshTokensResponse(resp)
	if err != nil {
		return res, errors.Wrap(err, "decode response")
	}

	return result, nil
}
