package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL        string
	AppPort            string
	Env                string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	RedisHost          string
	RedisPort          string
	RedisPassword      string
	JWTAccessSecret    string
	JWTRefreshSecret   string
	JWTAccessTTL       time.Duration
	JWTRefreshTTL      time.Duration
	FreeDailyPlayLimit int
	FreePlaylistLimit  int
}

func Load() *Config {
	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		AppPort:            getEnv("APP_PORT", "8080"),
		Env:                getEnv("ENV", "development"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "music_user"),
		DBPassword:         getEnv("DB_PASSWORD", "music_password"),
		DBName:             getEnv("DB_NAME", "music_service"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnv("REDIS_PORT", "6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		JWTAccessSecret:    getEnv("JWT_ACCESS_SECRET", "super_access_secret_key_change_me"),
		JWTRefreshSecret:   getEnv("JWT_REFRESH_SECRET", "super_refresh_secret_key_change_me"),
		JWTAccessTTL:       getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute),
		JWTRefreshTTL:      getEnvDuration("JWT_REFRESH_TTL", 720*time.Hour),
		FreeDailyPlayLimit: getEnvInt("FREE_DAILY_PLAY_LIMIT", 10),
		FreePlaylistLimit:  getEnvInt("FREE_PLAYLIST_LIMIT", 3),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if val, err := strconv.Atoi(value); err == nil {
			return val
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if val, err := time.ParseDuration(value); err == nil {
			return val
		}
	}
	return fallback
}
