package token_test

import (
	"context"
	"net"
	"testing"

	_ "github.com/lib/pq"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/aspirin100/JWT-API/internal/token"
)

func TestService_RefreshTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	keys := map[string]string{"testKeyID": "someTestSecretKey"}
	repo := NewMockRepository(ctrl)
	notifier := NewMockNotifier(ctrl)

	service := token.Service{
		SecretKeys:         keys,
		CurrentSecretKeyID: "testKeyID",
		RefreshTokenTTL:    43200,
		AccessTokenTTL:     15,
		Repository:         repo,
		Notifier:           notifier,
	}

	repo.EXPECT().
		InsertRefreshToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()

	userID := uuid.MustParse("a0fbc244-b17a-4a96-af2b-a3d7c058dc9f")

	accessToken, refreshToken, err := service.CreateNewTokensPair(context.Background(), &token.PairParams{
		IP:     net.ParseIP("1.2.3.4"),
		UserID: userID,
		AccessTokenOptClaims: jwt.MapClaims{
			"key1": "val1",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, accessToken)
	require.NotNil(t, refreshToken)

	repo.EXPECT().
		BeginTx(gomock.Any()).
		Return(context.Background(), token.CommitOrRollback(func(_ *error) error { return nil }), nil)
	repo.EXPECT().
		SetRefreshTokenUsed(gomock.Any(), gomock.Any()).
		Return(nil)

	accessToken, refreshToken, err = service.RefreshTokenPair(context.Background(), &token.RefreshTokensParams{
		IP: net.ParseIP("1.2.3.4"),
		AccessTokenOptClaims: jwt.MapClaims{
			"key1": "val1",
		},
		UserID:       userID,
		RefreshToken: *refreshToken,
		AccessToken:  *accessToken,
	})
	require.NoError(t, err)
	require.NotNil(t, accessToken)
	require.NotNil(t, refreshToken)
}
