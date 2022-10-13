package gredis

import (
	"context"
	"fmt"
	"log"

	gredis "github.com/go-redis/redis/v8"
)

var conn = gredis.NewClient(&gredis.Options{
	Network: "tcp",
	Addr:    "localhost:6379",
})

func init() {
	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
}

func Set(key, value string) error {

	ctx := context.Background()

	_, err := conn.Set(ctx, key, value, 0).Result()
	if err != nil {
		return fmt.Errorf("error while doing SET command in gredis : %v", err)
	}

	return err

}

func Get(key string) (string, error) {

	ctx := context.Background()

	value, err := conn.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("error while doing GET command in gredis : %v", err)
	}

	return value, err
}

func Hset(hash, key, value string) error {

	ctx := context.Background()

	_, err := conn.HSet(ctx, hash, key, value).Result()
	if err != nil {
		return fmt.Errorf("error while doing HSET command in gredis : %v", err)
	}

	return err
}

func Hget(hash, key string) (string, error) {

	ctx := context.Background()

	value, err := conn.HGet(ctx, hash, key).Result()
	if err != nil {
		return value, fmt.Errorf("error while doing HGET command in gredis : %v", err)
	}

	return value, err

}

func Hgetall(hash string) (map[string]string, error) {
	ctx := context.Background()

	value, err := conn.HGetAll(ctx, hash).Result()
	if err != nil {
		return value, fmt.Errorf("error while doing HGETALL command in gredis : %v", err)
	}

	return value, err
}
