package redis

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type Service struct {
	client *redis.Client
}

var service *Service = &Service{}
var redisServiceLock sync.Once

func GetService() *Service {
	if (*service).client == nil {
		redisServiceLock.Do(func() {
			client := redis.NewClient(&redis.Options{
				Addr:     "127.0.0.1:6379",
				Password: "",
				DB:       10,
			})

			_, err := client.Ping().Result()
			if err != nil {
				log.Fatal(err)
			}
			service = &Service{
				client: client,
			}
		})
	}
	return service
}

func (srv *Service) Get(key string) (interface{}, error) {
	val, err := srv.client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (srv *Service) MGet(keys []string) ([]interface{}, error) {
	return srv.client.MGet(keys...).Result()
}

func (srv *Service) Set(key string, value interface{}) error {
	exp := time.Duration(600 * time.Second)
	return srv.client.Set(key, value, exp).Err()
}

func (srv *Service) MSet(data map[string]interface{}) error {
	exp := time.Duration(600 * time.Second)
	var interfaces []interface{}
	pipe := srv.client.TxPipeline()
	for key, value := range data {
		interfaces = append(interfaces, key, value)
		pipe.Expire(key, exp)
	}

	if err := srv.client.MSet(interfaces...).Err(); err != nil {
		return err
	}
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}
