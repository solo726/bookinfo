package global

import (
	"log"
	"github.com/go-redis/redis"
)

func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     Conf.Redis.Addr,
		Password: Conf.Redis.Pwd,
		DB:       Conf.Redis.DB,
		//MaxRetries:   conf.MaxRetries,
		//PoolSize:     conf.PoolSize,
		//DialTimeout:  conf.DialTimeout,
		//ReadTimeout:  conf.ReadTimeout,
		//WriteTimeout: conf.WriteTimeout,
		//PoolTimeout:  conf.PoolTimeout,
		//IdleTimeout:  conf.IdleTimeout,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalln(pong, err)
	}
	return client
}
