package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
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

func IncrementBy(key string, value int64) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("client Redis not initialized")
	}
	newValue, err := RedisClient.IncrBy(ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementing key '%s' by %d: %w", key, value, err)
	}
	return newValue, nil
}

func Increment(key string) (int64, error) {
	return IncrementBy(key, 1)
}

func DecrementBy(key string, value int64) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("client Redis not initialized")
	}
	newValue, err := RedisClient.DecrBy(ctx, key, value).Result()
	if err != nil {
		return 0, fmt.Errorf("error decrementing key '%s' by %d: %w", key, value, err)
	}
	return newValue, nil
}

func Decrement(key string) (int64, error) {
	return DecrementBy(key, 1)
}

func GetInt(key string) (int64, error) {
	if RedisClient == nil {
		return 0, fmt.Errorf("client Redis not initialized")
	}

	valStr, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, redis.Nil
	} else if err != nil {
		return 0, fmt.Errorf("error getting key '%s' from redis: %w", key, err)
	}

	valInt, errConv := strconv.ParseInt(valStr, 10, 64)
	if errConv != nil {
		return 0, fmt.Errorf("value for key '%s' ('%s') could not be parsed as int64: %w", key, valStr, errConv)
	}

	return valInt, nil
}
