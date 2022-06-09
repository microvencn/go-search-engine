package database

import (
	"github.com/gomodule/redigo/redis"
	"go-search-engine/src/config"
	"log"
	"time"
)

var RedisClient *redis.Pool

func init() {
	redisConf := config.GetRedisConfig()
	RedisClient = &redis.Pool{
		MaxIdle:     redisConf.MaxIdle,
		MaxActive:   redisConf.MaxActive,
		IdleTimeout: time.Duration(redisConf.TimeOut) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(redisConf.Type, redisConf.Redis_Host)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			if _, err := c.Do("AUTH", redisConf.AUTH); err != nil {
				_ = c.Close()
				log.Println(err.Error())
				return nil, err
			}
			return c, nil
		},
	}
}
