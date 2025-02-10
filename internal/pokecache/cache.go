package pokecache

import (
	"encoding/gob"
	"os"
	"sync"
	"time"
)

type cacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

type Cache struct {
	cacheMap map[string]cacheEntry
	mu       *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cacheMap: make(map[string]cacheEntry),
		mu:       &sync.Mutex{},
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheMap[key] = cacheEntry{
		CreatedAt: time.Now().UTC(),
		Val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.cacheMap[key]
	return val.Val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *Cache) reap(now time.Time, interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.cacheMap {
		if v.CreatedAt.Before(now.Add(-interval)) {
			delete(c.cacheMap, k)
		}
	}
}

func (c *Cache) SaveToFile(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(c.cacheMap)
}

func (c *Cache) LoadFromFile(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&c.cacheMap)
}
