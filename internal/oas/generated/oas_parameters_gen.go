// Code generated by ogen, DO NOT EDIT.

package gen

import (
	"net/http"
	"net/url"

	"github.com/go-faster/errors"
	"github.com/google/uuid"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

// CreateTokensParams is parameters of CreateTokens operation.
type CreateTokensParams struct {
	UserID uuid.UUID
}

func unpackCreateTokensParams(packed middleware.Parameters) (params CreateTokensParams) {
	{
		key := middleware.ParameterKey{
			Name: "user_id",
			In:   "path",
		}
		params.UserID = packed[key].(uuid.UUID)
	}
	return params
}

func decodeCreateTokensParams(args [1]string, argsEscaped bool, r *http.Request) (params CreateTokensParams, _ error) {
	// Decode path: user_id.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "user_id",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToUUID(val)
				if err != nil {
					return err
				}

				params.UserID = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "user_id",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}

// RefreshTokensParams is parameters of RefreshTokens operation.
type RefreshTokensParams struct {
	UserID uuid.UUID
}

func unpackRefreshTokensParams(packed middleware.Parameters) (params RefreshTokensParams) {
	{
		key := middleware.ParameterKey{
			Name: "user_id",
			In:   "path",
		}
		params.UserID = packed[key].(uuid.UUID)
	}
	return params
}

func decodeRefreshTokensParams(args [1]string, argsEscaped bool, r *http.Request) (params RefreshTokensParams, _ error) {
	// Decode path: user_id.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "user_id",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToUUID(val)
				if err != nil {
					return err
				}

				params.UserID = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "user_id",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}
