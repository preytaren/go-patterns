package config

import "testing"

func TestNewConfig(t *testing.T) {
	cfg := NewConfig(1, WithTimeout(), WithC(10))

	if cfg.A != 1 {
		t.Error("cfg.A not equal to 1, actual %d", cfg.A)
	}

	if cfg.opts.Timeout != true {
		t.Error("cfg.timeout not true, actual %b", cfg.opts.Timeout)
	}

	if cfg.opts.C != 10 {
		t.Error("cfg.C not equal to 10, actual %d", cfg.opts.C)
	}
}