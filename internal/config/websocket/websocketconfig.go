package websocketconfig

import "time"

type WebSocketConfig struct {
	PongWait       time.Duration `yaml:"pong_wait"`
	PingPeriod     time.Duration `yaml:"ping_period"`
	WriteWait      time.Duration `yaml:"write_wait"`
	MaxMessageSize int64         `yaml:"max_message_size"`
}
