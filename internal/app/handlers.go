package app

import (
	"context"
	"net/http"

	"github.com/Vlad-Ali/Virtual-chat/internal/adapter/chat"
	"github.com/Vlad-Ali/Virtual-chat/internal/domain/session"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type Handlers struct {
	chatHandler *chat.Handler
	upgrader    *websocket.Upgrader
}

func NewHandlers(cfg *Config, service session.Service, appCtx context.Context) (*Handlers, http.Handler) {
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.CorsConfig.AllowedOrigins,
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}

			for _, allowedOrigin := range cfg.CorsConfig.AllowedOrigins {
				if allowedOrigin == origin {
					return true
				}
			}
			return false
		},
	}

	chatHandler := chat.NewHandler(&cfg.WebSocketConfig, upgrader, service)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/chat", chatHandler.ReceiveMsg(appCtx))

	return &Handlers{chatHandler, upgrader}, c.Handler(mux)
}
