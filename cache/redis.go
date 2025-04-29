package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var RedisClient *redis.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis connected")
}

func Set(key string, value string) error {
	if RedisClient == nil {
		return fmt.Errorf("client Redis not initialized")
	}
	return RedisClient.Set(ctx, key, value, 0).Err()
}

func SetStruct(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error serializing struct: %w", err)
	}
	return RedisClient.Set(ctx, key, jsonData, 0).Err()
}

func Get(key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("client Redis not initialized")
	}

	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}

func GetStruct(key string, dest interface{}) error {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error when fetching value: %w", err)
	}
	return json.Unmarshal([]byte(val), dest)
}

func Exists(key string) (bool, error) {
	if RedisClient == nil {
		return false, fmt.Errorf("client Redis not initialized")
	}
	val, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência da chave '%s': %w", key, err)
	}
	return val == 1, nil
}

func SetTTL(key string, expiration time.Duration) (bool, error) {
	if RedisClient == nil {
		return false, fmt.Errorf("client Redis not initialized")
	}
	return RedisClient.Expire(ctx, key, expiration).Result()
}

func Delete(key string) error {
	if RedisClient == nil {
		return fmt.Errorf("cliente Redis não inicializado")
	}
	return RedisClient.Del(ctx, key).Err()
}
