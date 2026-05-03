package session

import (
	"context"

	"github.com/Vlad-Ali/Virtual-chat/internal/domain/message"
)

type Provider interface {
	GetAnswer(ctx context.Context, history []*message.Message) (string, error)
}
