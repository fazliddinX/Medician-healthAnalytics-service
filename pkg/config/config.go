package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"

	_ "github.com/lib/pq"
)

type Config struct {
	HOST             string
	GRPC_SERVER_PORT string
	MONGO_URL        string
	RABBIT_URL       string
}

func Load() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	config := Config{}

	config.HOST = cast.ToString(coalesce("GIN_SERVER_PORT", ":8081"))
	config.GRPC_SERVER_PORT = cast.ToString(coalesce("GRPC_SERVER_PORT", ":50050"))
	config.MONGO_URL = cast.ToString(coalesce("MONGO_URL", "mongodb://mongo:27017"))
	config.RABBIT_URL = cast.ToString(coalesce("RABBIT_URL", "amqp://guest:guest@rabbitmq:5672/"))

	return config
}

func coalesce(env string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(env)
	if !exists {
		return defaultValue
	}
	return value
}
