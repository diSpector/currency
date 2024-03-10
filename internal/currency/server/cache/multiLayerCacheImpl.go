package cache

import (
	"time"

	"github.com/diSpector/currency.git/pkg/currency/entities"
	"github.com/pkg/errors"
)

type MultiLayerCache struct {
	head *cacheLayer
}

type cacheLayer struct {
	cache Cache
	next  *cacheLayer
}

func NewMultiLayerCache(cls ...Cache) Cache {
	cs := &MultiLayerCache{
		head: &cacheLayer{
			cache: cls[0],
		},
	}

	cur := cs.head
	if len(cls) > 1 {
		for i := 1; i < len(cls); i++ {
			cur.next = &cacheLayer{
				cache: cls[i],
			}
			cur = cur.next
		}
	}

	return cs
}

func (s *MultiLayerCache) Get(key string) (*entities.Currency, error) {
	cur := s.head
	for cur != nil {
		val, err := cur.cache.Get(key)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				cur = cur.next
				continue
			} else if errors.Is(err, ErrFoundInAbsent) {
				return nil, ErrFoundInAbsent
			} else {
				return nil, err
			}
		}
		return val, nil
	}

	return nil, ErrNotFound
}

func (s *MultiLayerCache) Set(key string, v *entities.Currency, ttl time.Duration) error {
	cur := s.head
	for cur != nil {
		_, err := cur.cache.Get(key)
		if err != nil {
			if errors.Is(err, ErrNotFound) || errors.Is(err, ErrFoundInAbsent) {
				errSet := cur.cache.Set(key, v, ttl)
				if errSet != nil {
					return errSet
				}
			} else {
				return err
			}
		}
		cur = cur.next
	}

	return nil
}

func (s *MultiLayerCache) Delete(k string) error {
	return nil
}

func (s *MultiLayerCache) List() error {
	return nil
}
