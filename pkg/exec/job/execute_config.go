package job

import "time"

type ExecuteConfig interface {
	isRunConfig()
}

type IntervalConfig struct {
	Interval time.Duration
	Timeout  time.Duration
	Retries  int
}

func (IntervalConfig) isRunConfig() {}

type ProcessConfig struct{}

func (ProcessConfig) isRunConfig() {}
