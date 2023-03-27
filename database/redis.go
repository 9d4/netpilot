package database

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var redisCli *redis.Client

func ConnectRedis(config *RedisConfig) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	redisCli = cli
	return cli, nil
}

func RedisCli() *redis.Client {
	return redisCli
}
