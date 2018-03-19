package linenotify

import (
	"sync"
	"time"

	"github.com/utahta/momoclo-channel/timeutil"
)

type (
	cache struct {
		buff      []byte
		expiredAt time.Time
	}

	cacheRepository struct {
		mux   sync.RWMutex
		cache map[string]*cache
	}
)

const (
	cacheExpiredInterval = 300 // sec
	cacheVacuumInterval  = 600 // sec
)

func newCache(b []byte) *cache {
	return &cache{buff: b, expiredAt: timeutil.Now().Add(cacheExpiredInterval * time.Second)}
}

func (c *cache) Bytes() []byte {
	return c.buff
}

func (c *cache) Expired() bool {
	return timeutil.Now().Unix() >= c.expiredAt.Unix()
}

func newCacheRepository() *cacheRepository {
	repo := &cacheRepository{cache: map[string]*cache{}}
	go func() {
		tick := time.NewTicker(cacheVacuumInterval * time.Second)
		for {
			select {
			case <-tick.C:
				repo.Vacuum()
			}
		}
	}()
	return repo
}

func (repo *cacheRepository) Set(k string, c *cache) {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	repo.cache[k] = c
}

func (repo *cacheRepository) Get(k string) (*cache, bool) {
	repo.mux.RLock()
	defer repo.mux.RUnlock()

	v, ok := repo.cache[k]
	return v, ok
}

func (repo *cacheRepository) Vacuum() {
	repo.mux.Lock()
	defer repo.mux.Unlock()

	for k, v := range repo.cache {
		if v.Expired() {
			delete(repo.cache, k)
		}
	}
}
