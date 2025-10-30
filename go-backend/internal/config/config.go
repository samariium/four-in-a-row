package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port             string
	MongoURI         string
	MatchBotAfterMs  int
	RejoinGraceMs    int
	BotMoveDelayMs   int
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func geti(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func Load() Config {
	return Config{
		Port:            getenv("PORT", "9090"),
		MongoURI:        getenv("MONGO_URI", "mongodb://localhost:27017"),
		MatchBotAfterMs: geti("MATCH_BOT_AFTER_MS", 10000),
		RejoinGraceMs:   geti("REJOIN_GRACE_MS", 30000),
		BotMoveDelayMs:  geti("BOT_MOVE_DELAY_MS", 400),
	}
}
