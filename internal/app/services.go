package app

import (
	"github.com/Vlad-Ali/Virtual-chat/internal/domain/session"
	"github.com/Vlad-Ali/Virtual-chat/internal/usecase"
)

type Services struct {
	Service session.Service
}

func NewServices(provider session.Provider) *Services {
	return &Services{Service: usecase.NewSessionService(provider)}
}
