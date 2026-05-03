package httpconfig

import "time"

type HTTPConfig struct {
	Address      string
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}
