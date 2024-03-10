package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/diSpector/currency.git/pkg/currency/entities"
	"github.com/redis/go-redis/v9"
)

const (
	NAMESPACE_ITEMS  = `items`
	NAMESPACE_ABSENT = `absent`
)

type RedisCache struct {
	client  *redis.Client
	ctxTime time.Duration
}

func NewRedisCache(client *redis.Client, t time.Duration) *RedisCache {
	return &RedisCache{
		client:  client,
		ctxTime: t,
	}
}

func (s *RedisCache) Get(key string) (*entities.Currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTime)
	defer cancel()

	// items:USD
	redisKey := NAMESPACE_ITEMS + `:` + key

	res, err := s.client.Get(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			// absent:NON
			absentKey := NAMESPACE_ABSENT + `:` + key
			_, err := s.client.Get(ctx, absentKey).Result()
			if err != nil {
				if err == redis.Nil {
					return nil, ErrNotFound
				} else {
					return nil, err
				}
			}

			log.Println(`Got from Redis absent:`, key)
			return nil, ErrFoundInAbsent

		} else {
			return nil, err
		}
	}

	var cur entities.Currency
	err = json.Unmarshal([]byte(res), &cur)
	if err != nil {
		return nil, err
	}

	log.Println(`Got from Redis items:`, key)

	return &cur, nil
}

func (s *RedisCache) Set(key string, val *entities.Currency, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ctxTime)
	defer cancel()

	var redisKey string
	if val != nil {
		redisKey = NAMESPACE_ITEMS + `:` + key
		log.Println(`Set to Redis items:`, key)
	} else {
		redisKey = NAMESPACE_ABSENT + `:` + key
		log.Println(`Set to Redis absent:`, key)
	}

	valByte, err := json.Marshal(val)
	if err != nil {
		return err
	}

	_, err = s.client.Set(ctx, redisKey, valByte, ttl).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *RedisCache) Delete(key string) error {
	return nil
}

func (s *RedisCache) List() error {
	// items:*
	// items:USD
	// absent:*
	return nil
}
