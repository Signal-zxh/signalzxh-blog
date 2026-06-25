package db

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() error {

	dbNum, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return err
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNum,
	})
	return RDB.Ping(context.Background()).Err()
}
