package global

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"golang.org/x/net/context"
	"github.com/openzipkin/zipkin-go"
	"github.com/pquerna/ffjson/ffjson"
)

type redisClient struct {
	client *redis.Client
}

func newRedisClient() *redisClient {
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
	return &redisClient{
		client: client,
	}
}

func (this *redisClient) WarpGet(ctx context.Context, key string) *redis.StringCmd {
	span, _ := ZPTracer.StartSpanFromContext(
		ctx,
		"getting-redis-info",
		zipkin.Tags(map[string]string{"func": "warpGet"}),
	)
	span.Annotate(time.Now(), "detail warpGet--start...")
	v := this.client.Get(key)
	defer func() {
		span.Annotate(time.Now(), "detail warpGet--end...")
		span.Finish()
	}()
	return v
}

func (this *redisClient) WarpSet(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration) *redis.StatusCmd {

	span, _ := ZPTracer.StartSpanFromContext(
		ctx,
		"setting-redis-info",
		zipkin.Tags(map[string]string{"func": "warpSet"}),
	)

	span.Annotate(time.Now(), "detail warpSet--start...")

	bytes, _ := ffjson.Marshal(value)

	defer func() {
		span.Annotate(time.Now(), "detail warpSet--end...")
		span.Finish()
	}()

	return  this.client.Set(key, bytes, expiration)

}