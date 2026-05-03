package chat

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	websocketconfig "github.com/Vlad-Ali/Virtual-chat/internal/config/websocket"
	"github.com/Vlad-Ali/Virtual-chat/internal/domain/message"
	"github.com/Vlad-Ali/Virtual-chat/internal/domain/session"
	"github.com/Vlad-Ali/Virtual-chat/internal/dto"
	"github.com/gorilla/websocket"
)

type Handler struct {
	config   *websocketconfig.WebSocketConfig
	logger   *slog.Logger
	upgrader *websocket.Upgrader
	service  session.Service
}

func NewHandler(cfg *websocketconfig.WebSocketConfig, upgrader *websocket.Upgrader, service session.Service) *Handler {
	return &Handler{config: cfg, logger: slog.Default().With("component", "chat_handler"), upgrader: upgrader, service: service}
}

func (h *Handler) ReceiveMsg(appCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		conn, err := h.upgrader.Upgrade(w, req, nil)
		if err != nil {
			slog.Error("error upgrading websocket connection", "error", err)
			return
		}
		defer conn.Close()
		ctx, cancel := context.WithCancel(appCtx)

		conn.SetReadLimit(h.config.MaxMessageSize)
		conn.SetReadDeadline(time.Now().Add(h.config.PongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(h.config.PongWait))
			return nil
		})
		inChan := make(chan string, 16)
		outChan := make(chan *message.Message, 16)

		errChan := h.service.HandleSession(ctx, inChan, outChan)

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			defer close(inChan)
			for {
				_, data, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
						h.logger.Debug("client closed connection gracefully")
					} else {
						h.logger.Error("read error", "error", err)
					}
					cancel()
					return
				}

				var msgObj dto.Message
				if jsonErr := json.Unmarshal(data, &msgObj); jsonErr != nil {
					h.logger.Error("error unmarshalling websocket message", "error", jsonErr)
					conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, "JSON Invalid format"), time.Now().Add(h.config.WriteWait))
					cancel()
					return
				}
				select {
				case inChan <- msgObj.Msg:
				case <-ctx.Done():
					return
				}
			}
		}()

		go func() {
			defer wg.Done()
			defer close(outChan)

			pingTicker := time.NewTicker(h.config.PingPeriod)
			defer pingTicker.Stop()

			for {
				select {
				case msg, ok := <-outChan:
					conn.SetWriteDeadline(time.Now().Add(h.config.WriteWait))
					if !ok {
						conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(h.config.WriteWait))
						cancel()
						return
					}

					var response dto.Message
					response.Msg = msg.Msg

					jsonData, jsonErr := json.Marshal(response)
					if jsonErr != nil {
						h.logger.Error("error marshalling websocket message", "error", jsonErr)
						conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, ""), time.Now().Add(h.config.WriteWait))
						cancel()
						return
					}

					conn.SetWriteDeadline(time.Now().Add(h.config.WriteWait))
					if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
						h.logger.Error("error writing websocket message", "error", err)
						cancel()
						return
					}

				case <-pingTicker.C:
					conn.SetWriteDeadline(time.Now().Add(h.config.WriteWait))
					if wrErr := conn.WriteMessage(websocket.PingMessage, []byte("ping")); wrErr != nil {
						h.logger.Error("error pinging websocket", "error", wrErr)
						cancel()
						return
					}
				case err, ok := <-errChan:
					if !ok {
						errChan = nil
						continue
					}

					if err != nil {
						h.logger.Error("error handling response message", "error", err)
						if errors.Is(err, message.ErrMessageInvalidFormat) {
							conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, "Message invalid format"), time.Now().Add(h.config.WriteWait))
						} else {
							conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, ""), time.Now().Add(h.config.WriteWait))
						}
						cancel()
						return
					}
				case <-ctx.Done():
					conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(h.config.WriteWait))
					h.logger.Debug("closed connection gracefully", "error", err)
					return
				}
			}
		}()

		wg.Wait()

		h.logger.Info("websocket connection fully closed")
	}
}
