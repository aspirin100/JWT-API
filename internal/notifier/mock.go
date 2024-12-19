package notifier

import (
	"context"

	"github.com/google/uuid"

	"github.com/aspirin100/JWT-API/internal/logger"
)

type Mock struct{}

func (m *Mock) Notify(_ context.Context, userID uuid.UUID, subject, message string) error {
	logger.Default().Info("mock notify", "userID", userID, "subject", subject, "message", message)

	return nil
}
