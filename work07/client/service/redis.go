package service

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log"
)

var RDB *redis.Client
var err error

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "119.3.125.7:6379",
		DB:   15,
	})
}

func RdbLPush(key string, values interface{}) error {
	_, err = RDB.LPush(context.Background(), key, values).Result()
	if err != nil {
		return err
	}
	return nil
}

func RdbLLen(key string) (int64, error) {
	res, err := RDB.LLen(context.Background(), key).Result()
	if err != nil {
		log.Println("RDB.LLen failed")
	}
	return res, nil
}

func RdbLPop(key string) (string, error) {
	res, err := RDB.LPop(context.Background(), key).Result()

	if err != nil {
		log.Println("RDB.LPop failed")
	}
	return res, nil
}
