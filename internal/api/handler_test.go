package api_test

import (
	"context"
	"net"
	"testing"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/aspirin100/JWT-API/internal/api"
	"github.com/aspirin100/JWT-API/internal/middleware"
	gen "github.com/aspirin100/JWT-API/internal/oas/generated"
	"github.com/aspirin100/JWT-API/internal/token"
)

func TestHandler_CreateTokensSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tokenService := NewMockTokenService(ctrl)
	handler := api.Handler{
		TokenService: tokenService,
	}

	ip := net.ParseIP("127.0.0.1")
	ctx := middleware.SetClientIP(context.Background(), ip)
	accessToken := "access_token"
	refreshToken := "refresh_token"

	tokenService.EXPECT().
		CreateNewTokensPair(ctx, &token.PairParams{
			IP:     ip,
			UserID: uuid.MustParse("9bdc7259-fe94-4829-bcbc-6f1be24a7c34"),
			AccessTokenOptClaims: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		}).
		Return(
			&accessToken,
			&refreshToken,
			nil,
		)

	res, err := handler.CreateTokens(ctx, gen.OptCreateTokensRequest{
		Value: gen.CreateTokensRequest{
			AdditionalInfo: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		Set: true,
	}, gen.CreateTokensParams{
		UserID: uuid.MustParse("9bdc7259-fe94-4829-bcbc-6f1be24a7c34"),
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	resp, ok := res.(*gen.CreateTokensResponse)
	require.True(t, ok)
	require.NotEmpty(t, resp.AccessToken)
	require.NotEmpty(t, resp.RefreshToken)
}
