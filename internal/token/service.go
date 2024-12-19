package token

import (
	"context"
	"fmt"
	"time"

	"github.com/aspirin100/JWT-API/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CommitOrRollback func(err *error) error

//go:generate go run go.uber.org/mock/mockgen@latest -source=service.go -destination=service_mocks_test.go -package=token_test

type Notifier interface {
	Notify(ctx context.Context, userID uuid.UUID, subject, message string) error
}

type Repository interface {
	BeginTx(ctx context.Context) (context.Context, CommitOrRollback, error)
	InsertRefreshToken(ctx context.Context, pairID, userID uuid.UUID, refreshToken string) error
	SetRefreshTokenUsed(ctx context.Context, pairID uuid.UUID) error
}

type Service struct {
	SecretKeys         map[string]string
	CurrentSecretKeyID string
	RefreshTokenTTL    time.Duration
	AccessTokenTTL     time.Duration

	Notifier   Notifier
	Repository Repository
}

const (
	pairIDClaimKey    = "pairID"
	ipClaimKey        = "ip"
	expiresAtClaimKey = "expiresAt"
	userIDClaimKey    = "userID"
	claimsClaimKey    = "claims"
)

func (s *Service) CreateNewTokensPair(ctx context.Context, params *PairParams) (_, _ *string, err error) {
	pairID := uuid.New()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		pairIDClaimKey:    pairID,
		ipClaimKey:        params.IP,
		expiresAtClaimKey: time.Now().Add(s.RefreshTokenTTL).Format(time.RFC3339),
	})

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		pairIDClaimKey:    pairID,
		ipClaimKey:        params.IP,
		userIDClaimKey:    params.UserID,
		claimsClaimKey:    params.AccessTokenOptClaims,
		expiresAtClaimKey: time.Now().Add(s.AccessTokenTTL).Format(time.RFC3339),
	})

	accessToken.Header["kid"] = s.CurrentSecretKeyID
	refreshToken.Header["kid"] = s.CurrentSecretKeyID

	refreshSigned, err := refreshToken.SignedString([]byte(s.SecretKeys[s.CurrentSecretKeyID]))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	accessSigned, err := accessToken.SignedString([]byte(s.SecretKeys[s.CurrentSecretKeyID]))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	err = s.Repository.InsertRefreshToken(ctx, pairID, params.UserID, refreshSigned)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to insert refresh token into database: %w", err)
	}

	logger.Default().Info("token pair was successfully created")

	return &accessSigned, &refreshSigned, nil
}

func (s *Service) RefreshTokenPair( //nolint:gocyclo,cyclop
	ctx context.Context,
	params *RefreshTokensParams,
) (_, _ *string, err error) {
	// проверка валидности refresh токена по времени
	refreshClaims := make(jwt.MapClaims)

	_, err = jwt.ParseWithClaims(params.RefreshToken, refreshClaims, s.keyFunc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse jwt with claims: %w", err)
	}

	expirationTime, err := time.Parse(time.RFC3339, refreshClaims[expiresAtClaimKey].(string))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse expiration time: %w", err)
	}

	if expirationTime.Unix() < time.Now().Unix() {
		return nil, nil, ErrRefreshTokenExpired
	}

	currentIP := params.IP
	currentUserID := params.UserID

	// проверка ip запроса на рефреш на соответствие с ip, содержащемся в старом refresh токене
	accessClaims := make(jwt.MapClaims)

	_, err = jwt.ParseWithClaims(params.AccessToken, accessClaims, s.keyFunc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse with claims: %w", err)
	}

	_, ok := accessClaims[userIDClaimKey].(string)
	if !ok {
		return nil, nil, ErrUserIDRequired
	}

	oldUserID, err := uuid.Parse(accessClaims[userIDClaimKey].(string))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse user id: %w", err)
	}

	if oldUserID != currentUserID {
		return nil, nil, fmt.Errorf("user id mismatch: %w", err)
	}

	_, ok = accessClaims[pairIDClaimKey].(string)
	if !ok {
		return nil, nil, ErrPairIDRequired
	}

	oldPairID, err := uuid.Parse(accessClaims[pairIDClaimKey].(string))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse user id from access token: %w", err)
	}

	oldIP, ok := accessClaims[ipClaimKey].(string)
	if !ok {
		return nil, nil, ErrIPRequired
	}

	if oldIP != currentIP.String() {
		// имитация отправки warning'a на email
		err = s.Notifier.Notify(ctx, oldUserID, "Warning", "IP changed from "+oldIP+" to "+currentIP.String())
		if err != nil {
			return nil, nil, fmt.Errorf("notify failure: %w", err)
		}
	}

	// начинаем транзакцию для добавления и обновления данных в postgres
	txCtx, commitOrRollback, err := s.Repository.BeginTx(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if errTx := commitOrRollback(&err); errTx != nil {
			err = errTx
		}
	}()

	err = s.Repository.SetRefreshTokenUsed(txCtx, oldPairID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set refresh token status used: %w", err)
	}

	var accessTokenOptClaims jwt.MapClaims

	if claims, ok := accessClaims[claimsClaimKey].(map[string]any); ok {
		accessTokenOptClaims = claims
	}

	accessToken, refreshToken, err := s.CreateNewTokensPair(txCtx, &PairParams{
		IP:                   currentIP,
		UserID:               currentUserID,
		AccessTokenOptClaims: accessTokenOptClaims,
	})
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, err
}

func (s *Service) keyFunc(token *jwt.Token) (any, error) {
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("empty kid: %w", ErrInvalidToken)
	}

	secret, ok := s.SecretKeys[kid]
	if !ok {
		return nil, fmt.Errorf("wrong secret key: %w", ErrInvalidToken)
	}

	return []byte(secret), nil
}
