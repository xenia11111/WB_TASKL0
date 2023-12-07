package cache

import (
	"time"

	app "github.com/xenia11111/WB_TASKL0"
	"github.com/xenia11111/WB_TASKL0/pkg/repository"
)

type OrderCRUD interface {
	Set(key string, value app.Order, duration time.Duration)

	Get(key string) (*app.OrderBody, bool)

	Delete(key string) error

	StartGC()

	GC()

	expiredKeys() (keys []string)

	clearOrderItems(keys []string)
}

type Cache struct {
	OrderCRUD
}

func NewCache(defaultExpiration, cleanupInterval time.Duration, repos repository.OrderCRUD, cacheItemLimit int) *Cache {
	return &Cache{
		OrderCRUD: NewOrderCRUD(defaultExpiration, cleanupInterval, repos, cacheItemLimit),
	}
}
