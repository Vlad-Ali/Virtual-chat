package session

import (
	"github.com/Vlad-Ali/Virtual-chat/internal/domain/message"
)

const (
	MaxHistorySize = 20
)

type Session struct {
	history []*message.Message
}

func NewSession() *Session {
	return &Session{history: []*message.Message{}}
}

func (session *Session) AddMessage(msg *message.Message) {
	if len(session.history) >= MaxHistorySize {
		for i := 1; i < len(session.history); i++ {
			session.history[i-1] = session.history[i]
		}
		session.history[MaxHistorySize-1] = msg
	} else {
		session.history = append(session.history, msg)
	}
}

func (session *Session) History() []*message.Message {
	result := make([]*message.Message, len(session.history))
	copy(result, session.history)
	return result
}
