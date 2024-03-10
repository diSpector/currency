package cache

import (
	"time"

	"github.com/diSpector/currency.git/pkg/currency/entities"
)

type Cache interface {
	Get(k string) (*entities.Currency, error)
	Set(k string, v *entities.Currency, ttl time.Duration) error
	Delete(k string) error
	List() error
}