package common

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	items sync.Map
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	expiration := time.Now().Add(duration).UnixNano()
	c.items.Store(key, CacheItem{
		Value:      value,
		Expiration: expiration,
	})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, found := c.items.Load(key)
	if !found {
		return nil, false
	}
	cacheItem := item.(CacheItem)
	if time.Now().UnixNano() > cacheItem.Expiration {
		c.items.Delete(key)
		return nil, false
	}
	return cacheItem.Value, true
}

func (c *Cache) Delete(key string) {
	c.items.Delete(key)
}

func (c *Cache) cleanupExpiredItems() {
	for {
		time.Sleep(time.Minute)
		now := time.Now().UnixNano()
		c.items.Range(func(key, value interface{}) bool {
			cacheItem := value.(CacheItem)
			if now > cacheItem.Expiration {
				c.items.Delete(key)
			}
			return true
		})
	}
}
func NewCache() *Cache {
	cache := &Cache{}
	go cache.cleanupExpiredItems()
	return cache
}

var MyCache = NewCache()
