package session

import (
	"context"

	"github.com/Vlad-Ali/Virtual-chat/internal/domain/message"
)

type Service interface {
	HandleSession(ctx context.Context, inChan <-chan string, outChan chan<- *message.Message) <-chan error
}
