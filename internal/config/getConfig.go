package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	RequestTimeout      time.Duration
	CalculateVolumeSize bool
}

func GetConfig() Config {
	timeout, err := time.ParseDuration(os.Getenv("REQUEST_TIMEOUT"))
	if err != nil || timeout < 0 {
		timeout = 120 * time.Second
	}

	envVolSize := strings.TrimSpace(strings.ToUpper(os.Getenv("CALCULATE_VOLUME_SIZE")))
	var volSize bool
	if envVolSize == "0" || envVolSize == "FALSE" {
		volSize = false
	} else {
		volSize = true
	}

	return Config{
		RequestTimeout:      timeout,
		CalculateVolumeSize: volSize,
	}
}
