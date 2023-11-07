package config

import (
	"os"
	"time"
)

type Config struct {
	RequestTimeout time.Duration
}

func GetConfig() Config {
	timeout, err := time.ParseDuration(os.Getenv("REQUEST_TIMEOUT"))
	if err != nil || timeout < 0 {
		timeout = 120 * time.Second
	}

	return Config{
		RequestTimeout: timeout,
	}
}
