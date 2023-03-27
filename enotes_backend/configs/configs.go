package configs

import "time"

func CurrTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

const (
	BEARER_SCHEMA = "Bearer "
)