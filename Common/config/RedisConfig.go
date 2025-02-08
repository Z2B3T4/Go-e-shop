package config

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func NewRedisClient(RedisDBNum int) *redis.Client {

	err := InitRedis(RedisDBNum)
	if err != nil {
		fmt.Println(err)
	}
	return redisClient
}

var redisClient *redis.Client

func InitRedis(RedisDBNum int) error {
	err := godotenv.Load("config.env")
	if err != nil {
		return nil
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"), // Redis地址和端口
		Password: os.Getenv("REDIS_PASSWORD"),                             // Redis密码
		DB:       RedisDBNum,                                              // 默认DB
	})
	_, err = redisClient.Ping().Result()
	return err
}

func SaveToken(userId int32, token string) error {
	return redisClient.Set(string(userId), token, 24*time.Minute).Err()
}

func GetToken(userId int32) (string, error) {
	return redisClient.Get(string(userId)).Result()
}
