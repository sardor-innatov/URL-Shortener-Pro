package config

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)


var (
	redisClient *redis.Client
	once        sync.Once
)
func GetRedis() *redis.Client {

	envProj := ProjectEnv()
	addr := envProj.RedisAddr
	password := envProj.RedisPassword

//	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     addr, 
			Password: password,
			DB:       0,
		})

		// Проверка соединения при старте
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			panic("Could not connect to Redis: " + err.Error())
		}
	//})
	return redisClient
}