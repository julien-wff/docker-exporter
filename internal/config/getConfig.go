package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	RequestTimeout time.Duration
}

func GetConfig() Config {
	timeout, err := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT"))
	if err != nil || timeout < 0 {
		timeout = 120
	}
	if timeout == 0 {
		timeout = 24 * 3600
	}

	return Config{
		RequestTimeout: time.Duration(timeout) * time.Second,
	}
}
