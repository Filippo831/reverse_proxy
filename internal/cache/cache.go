// implementing a caching system using the tutorial at this link https://mayurkhante786.medium.com/understanding-caching-in-go-part-1-improving-performance-and-efficiency-f9391e2d7047
// and adapting it to my usecase
package cache

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type Cache struct {
	data        map[*http.Request]http.ResponseWriter // store data
	expiration  map[*http.Request]time.Time           // store timeout for each key
	mutex       sync.Mutex
	defaultTTL  time.Duration // default time a data is stored into cache
	cleanupTick time.Duration // when this tick occurs delete expired data
	dataStored  int           // keep track of the amount of data stored to not exceed the size
	cacheSize   int
}

func NewCache(_defaultTTL time.Duration, _cleanupTick time.Duration) *Cache {
	cache := &Cache{
		data:        make(map[*http.Request]http.ResponseWriter),
		expiration:  make(map[*http.Request]time.Time),
		mutex:       sync.Mutex{},
		defaultTTL:  _defaultTTL,
		cleanupTick: _cleanupTick,
	}

	return cache
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.cleanupTick)
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		}
	}
}

func (c *Cache) cleanup() {
	currentTime := time.Now()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, el := range c.expiration {
		if currentTime.After(el) {
			delete(c.data, key)
			delete(c.expiration, key)
		}
	}
}

func (c *Cache) Set(key *http.Request, data *http.ResponseWriter) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = *data
	c.expiration[key] = time.Now().Add(c.defaultTTL)

	return nil
}

func (c *Cache) Get(key *http.Request) (*http.ResponseWriter, bool) {
	value, ok := c.data[key]
	return &value, ok
}

func (c *Cache) Print() {
    log.Print(len(c.data))
}
