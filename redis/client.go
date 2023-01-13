package redis

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
)

type Service struct {
	client *redis.ClusterClient
}

var service *Service = &Service{}
var ctx = context.Background()
var redisServiceLock sync.Once

func GetService() *Service {
	if (*service).client == nil {
		// redisServiceLock.Do(func() {
		// 	client := redis.NewClient(&redis.Options{
		// 		Addr:     os.Getenv("REDIS_HOST"),
		// 		Password: "",
		// 		DB:       0,
		// 	})

		// 	_, err := client.Ping().Result()
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	service = &Service{
		// 		client: client,
		// 	}
		// })
		//Create a new cluster client
		redisServiceLock.Do(func() {
			client := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: []string{"173.18.0.2:6379", "173.18.0.3:6379", "173.18.0.4:6379", "173.18.0.5:6379", "173.18.0.6:6379", "173.18.0.7:6379"},
			})
			err := client.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
				return shard.Ping(ctx).Err()
			})
			if err != nil {
				panic(err)
			}

			service = &Service{
				client: client,
			}
		})
	}
	return service
}

func (srv *Service) Get(key string) (interface{}, error) {
	val, err := srv.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (srv *Service) MGet(keys []string) ([]interface{}, error) {
	return srv.client.MGet(ctx, keys...).Result()
}

func (srv *Service) Set(key string, value interface{}) error {
	exp := time.Duration(600 * time.Second)
	return srv.client.Set(ctx, key, value, exp).Err()
}

func (srv *Service) MSet(data map[string]interface{}) error {
	exp := time.Duration(600 * time.Second)
	var interfaces []interface{}
	pipe := srv.client.TxPipeline()
	for key, value := range data {
		interfaces = append(interfaces, key, value)
		pipe.Expire(ctx, key, exp)
	}

	if err := srv.client.MSet(ctx, interfaces...).Err(); err != nil {
		return err
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}
