package cache

import (
	"errors"
	"sync"
	"time"

	app "github.com/xenia11111/WB_TASKL0"
	"github.com/xenia11111/WB_TASKL0/pkg/repository"
)

type OrderCache struct {
	sync.RWMutex
	orderItems        map[string]OrderItem
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	cacheItemLimit    int
}

type OrderItem struct {
	Value      app.OrderBody
	Expiration int64
	Created    time.Time
}

func NewOrderCRUD(defaultExpiration, cleanupInterval time.Duration,
	repos repository.OrderCRUD,
	cacheItemLimit int) *OrderCache {

	order_items := make(map[string]OrderItem)

	orders, err := repos.GetBulk(cacheItemLimit)

	if err == nil {

		for _, order := range *orders {
			order_items[order.Order_uid] = OrderItem{
				Value:      order.Order_body,
				Expiration: int64(defaultExpiration),
				Created:    time.Now(),
			}
		}
	}

	cache := OrderCache{
		orderItems:        order_items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		cacheItemLimit:    cacheItemLimit,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

func (c *OrderCache) Set(orderUID string, value app.Order, duration time.Duration) {

	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.orderItems[orderUID] = OrderItem{
		Value:      value.Order_body,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

func (c *OrderCache) Get(key string) (*app.OrderBody, bool) {

	c.RLock()
	defer c.RUnlock()

	item, found := c.orderItems[key]
	if !found {
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}

	return &item.Value, true
}

func (c *OrderCache) Delete(key string) error {

	c.Lock()
	defer c.Unlock()

	if _, found := c.orderItems[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.orderItems, key)

	return nil
}

func (c *OrderCache) StartGC() {
	go c.GC()
}

func (c *OrderCache) GC() {
	for {

		<-time.After(c.cleanupInterval)

		if c.orderItems == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearOrderItems(keys)
		}
	}
}

func (c *OrderCache) expiredKeys() (keys []string) {

	c.RLock()
	defer c.RUnlock()

	for k, i := range c.orderItems {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (c *OrderCache) clearOrderItems(keys []string) {

	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.orderItems, k)
	}
}
