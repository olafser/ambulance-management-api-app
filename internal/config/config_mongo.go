package config

import (
	"net/url"
	"os"
	"strconv"
	"time"
)

type MongoConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	Database     string
	Timeout      time.Duration
	AuthSource   string
	VehiclesColl string
	CountersColl string
}

func LoadMongoConfig() MongoConfig {
	timeoutSeconds, err := strconv.Atoi(getEnv("AMBULANCE_MANAGEMENT_API_MONGODB_TIMEOUT_SECONDS", "10"))
	if err != nil || timeoutSeconds <= 0 {
		timeoutSeconds = 10
	}

	return MongoConfig{
		Host:         getEnv("AMBULANCE_MANAGEMENT_API_MONGODB_HOST", "localhost"),
		Port:         getEnv("AMBULANCE_MANAGEMENT_API_MONGODB_PORT", "27017"),
		Username:     os.Getenv("AMBULANCE_MANAGEMENT_API_MONGODB_USERNAME"),
		Password:     os.Getenv("AMBULANCE_MANAGEMENT_API_MONGODB_PASSWORD"),
		Database:     "ambulance_management",
		Timeout:      time.Duration(timeoutSeconds) * time.Second,
		AuthSource:   "admin",
		VehiclesColl: "vehicles",
		CountersColl: "counters",
	}
}

func (c MongoConfig) URI() string {
	if c.Username == "" {
		return "mongodb://" + c.Host + ":" + c.Port
	}

	cred := url.UserPassword(c.Username, c.Password).String()
	return "mongodb://" + cred + "@" + c.Host + ":" + c.Port + "/?authSource=" + url.QueryEscape(c.AuthSource)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
