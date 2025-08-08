package di

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DB    DBConfig
	Kafka KafkaConfig
	HTTP  HTTPConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers  []string
	Topic    string
	GroupID  string
	MinBytes int
	MaxBytes int
}

type HTTPConfig struct {
	Host string
	Port string
}

func ReadConfig() Config {
	db := DBConfig{
		Host:     mustGetEnv("DB_HOST"),
		Port:     mustGetEnv("DB_PORT"),
		Username: mustGetEnv("DB_USERNAME"),
		Password: mustGetEnv("DB_PASSWORD"),
		DBName:   mustGetEnv("DB_NAME"),
		SSLMode:  mustGetEnv("SSL_MODE"),
	}

	kafka := KafkaConfig{
		Brokers:  parseBrokers(mustGetEnv("KAFKA_BROKERS")),
		Topic:    mustGetEnv("KAFKA_TOPIC"),
		GroupID:  mustGetEnv("KAFKA_GROUP_ID"),
		MinBytes: mustGetEnvInt("KAFKA_MIN_BYTES"),
		MaxBytes: mustGetEnvInt("KAFKA_MAX_BYTES"),
	}

	http := HTTPConfig{
		Host: mustGetEnv("HTTP_HOST"),
		Port: mustGetEnv("HTTP_PORT"),
	}

	return Config{
		DB:    db,
		Kafka: kafka,
		HTTP:  http,
	}
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("environment variable %s is not set", key)
	}
	return val
}

func mustGetEnvInt(key string) int {
	valStr := mustGetEnv(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("invalid int value for %s: %v", key, err)
	}
	return val
}

func parseBrokers(raw string) []string {
	parts := strings.Split(raw, ",")
	var brokers []string
	for _, b := range parts {
		b = strings.TrimSpace(b)
		if b != "" {
			brokers = append(brokers, b)
		}
	}
	if len(brokers) == 0 {
		log.Fatal("no valid Kafka brokers provided in KAFKA_BROKERS")
	}
	return brokers
}
