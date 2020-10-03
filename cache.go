package GeeCache

import (
	"GeeCache/lru"
	"sync"
)

//并发控制
type cache struct {
	mu         sync.Mutex //互斥锁
	lru        *lru.Cache
	cacheBytes int64
}

//添加
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

//查询
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
