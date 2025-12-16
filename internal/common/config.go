package common

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	GinMode         string
	OpenAIAPIKey    string
	MaxAudioSizeMB  int64
	GoogleCredsFile string
	GoogleSheetID   string
}

var AppConfig *Config

func LoadConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using env vars")
	}

	maxSize, err := strconv.ParseInt(getEnv("MAX_AUDIO_SIZE_MB", "10"), 10, 64)
	if err != nil {
		maxSize = 10
	}

	AppConfig = &Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),
		MaxAudioSizeMB:  maxSize,
		GoogleCredsFile: getEnv("GOOGLE_CREDENTIALS_FILE", ""),
		GoogleSheetID:   getEnv("GOOGLE_SHEET_ID", ""),
	}

	if AppConfig.OpenAIAPIKey == "" {
		return NewBadRequestError("OPENAI_API_KEY required")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
