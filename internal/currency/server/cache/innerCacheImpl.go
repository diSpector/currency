package cache

import (
	"log"
	"sync"
	"time"

	"github.com/diSpector/currency.git/pkg/currency/entities"
)

type InnerCache struct {
	mx     sync.RWMutex
	items  map[string]item
	absent map[string]time.Time
}

type item struct {
	v   *entities.Currency
	ttl time.Time
}

func NewInnerCache() *InnerCache {
	s := &InnerCache{
		items:  make(map[string]item),
		absent: make(map[string]time.Time),
	}

	go s.collect()
	go s.collectAbsents()

	return s
}

func (s *InnerCache) Get(key string) (*entities.Currency, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if v, ok := s.items[key]; ok {
		log.Println(`Get from inner cache:`, key)
		return v.v, nil
	} else {
		if _, ok := s.absent[key]; ok {
			log.Println(`Get from inner cache absent:`, key)
			return nil, ErrFoundInAbsent
		} else { // !ok - key NOT found in absent
			return nil, ErrNotFound
		}
	}
}

func (s *InnerCache) Set(key string, val *entities.Currency, ttl time.Duration) error {
	s.mx.Lock()

	if val != nil {
		s.items[key] = item{
			v:   val,
			ttl: time.Now().Add(ttl),
		}
		log.Println(`Set to inner cache:`, key)
	} else {
		s.absent[key] = time.Now().Add(ttl)
		log.Println(`Set to inner cache absent:`, key)
	}

	s.mx.Unlock()
	return nil
}

func (s *InnerCache) Delete(key string) error {
	s.mx.Lock()
	delete(s.items, key)
	s.mx.Unlock()
	return nil
}

func (s *InnerCache) List() error {
	s.mx.RLock()
	for k, item := range s.items {
		log.Println(k, item)
	}
	s.mx.RUnlock()
	return nil
}

func (s *InnerCache) collect() {
	for {
		s.mx.Lock()
		for k, v := range s.items {
			if time.Now().After(v.ttl) {
				delete(s.items, k)
				log.Println(`collect:`, k)
			}
		}
		s.mx.Unlock()
		time.Sleep(30 * time.Second)
	}
}

func (s *InnerCache) collectAbsents() {
	for {
		s.mx.Lock()
		for k, v := range s.absent {
			if v.After(time.Now()) {
				delete(s.absent, k)
				log.Println(`collect absent:`, k)
			}
		}
		s.mx.Unlock()
		time.Sleep(2 * time.Minute)
	}
}
