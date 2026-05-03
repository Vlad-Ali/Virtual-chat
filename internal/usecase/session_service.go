package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/Vlad-Ali/Virtual-chat/internal/domain/message"
	sessiondomain "github.com/Vlad-Ali/Virtual-chat/internal/domain/session"
)

const (
	BatchSize   = 5
	DrainPeriod = time.Second * 5
	IdlePeriod  = time.Second * 30
)

type SessionService struct {
	provider sessiondomain.Provider
	logger   *slog.Logger
}

func NewSessionService(provider sessiondomain.Provider) *SessionService {
	return &SessionService{logger: slog.Default().With("component", "session_service"), provider: provider}
}

func (s *SessionService) HandleSession(ctx context.Context, inChan <-chan string, outChan chan<- *message.Message) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		session := sessiondomain.NewSession()

		isProcessing := atomic.Bool{}

		pendingBatch := make([]*message.Message, 0, BatchSize)

		drainTicker := time.NewTicker(DrainPeriod)
		idleTicker := time.NewTicker(IdlePeriod)
		defer drainTicker.Stop()
		defer idleTicker.Stop()

		historyHandler := func() {
			s.logger.Debug("Handling history session")
			answer, err := s.provider.GetAnswer(ctx, session.History())
			msg := &message.Message{Type: message.ModelType, Msg: answer}
			isProcessing.Store(false)

			if err != nil {
				select {
				case errChan <- fmt.Errorf("error getting answer %w", err):
					s.logger.Error("error getting answer ", "error", err)
					return
				case <-ctx.Done():
					return
				}
			}

			session.AddMessage(msg)

			isProcessing.Store(false)

			select {
			case outChan <- msg:
				s.logger.Debug("Sending message to outChan")
				return
			case <-ctx.Done():
				return
			}
		}

		for {
			select {
			case text, ok := <-inChan:
				idleTicker.Reset(IdlePeriod)
				if !ok {
					return
				}

				if text == "" {
					select {
					case errChan <- message.ErrMessageInvalidFormat:
					case <-ctx.Done():
						return
					}
				}

				msg := &message.Message{Msg: text, Type: message.UserType}

				pendingBatch = append(pendingBatch, msg)
			case <-drainTicker.C:
				if isProcessing.Load() {
					idleTicker.Reset(IdlePeriod)
					continue
				}

				if len(pendingBatch) == 0 {
					continue
				}

				idleTicker.Reset(IdlePeriod)

				for _, msg := range pendingBatch {
					session.AddMessage(msg)
				}

				for i := range pendingBatch {
					pendingBatch[i] = nil
				}

				pendingBatch = pendingBatch[:0]

				isProcessing.Store(true)
				go historyHandler()

			case <-idleTicker.C:
				if isProcessing.Load() {
					idleTicker.Reset(IdlePeriod)
					continue
				}

				callBackMsg := &message.Message{Msg: "", Type: message.CallBackType}
				session.AddMessage(callBackMsg)

				s.logger.Debug("Sending CallBack Message to provider")
				isProcessing.Store(true)
				go historyHandler()

			case <-ctx.Done():
				return
			}
		}
	}()

	return errChan
}
