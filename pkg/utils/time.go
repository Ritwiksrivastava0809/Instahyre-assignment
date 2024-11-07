package utils

import (
	"log"
	"spam-search/pkg/config"
	"time"
)

func GetAccessTokenDuration() time.Duration {

	durationStr := config.GetAccessTokenDuration()
	if durationStr == "" {
		durationStr = "15m"
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {

		log.Printf("Failed to parse ACCESS_TOKEN_DURATION: %v", err)
		return 15 * time.Minute
	}
	return duration
}
