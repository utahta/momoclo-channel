package linenotify

import (
	"sync"
	"testing"
	"time"

	"github.com/utahta/momoclo-channel/lib/timeutil"
)

func TestCache_Bytes(t *testing.T) {
	c := newCache([]byte("test"))
	if string(c.Bytes()) != "test" {
		t.Errorf("Expected test, got %v", string(c.Bytes()))
	}
}

func TestCache_Expired(t *testing.T) {
	c := newCache([]byte("test"))

	if c.expiredAt.IsZero() {
		t.Fatal("expiredAt is zero")
	}

	if c.Expired() {
		t.Errorf("Expected not expired, but expired")
	}

	tmp := timeutil.Now
	timeutil.Now = func() time.Time {
		return tmp().Add(cacheExpiredInterval * time.Second)
	}
	defer func() {
		timeutil.Now = tmp
	}()

	if !c.Expired() {
		t.Errorf("Expected expired, but not expired")
	}
}

func TestCacheRepository_Set(t *testing.T) {
	repo := &cacheRepository{cache: map[string]*cache{}}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			repo.Set("a", newCache([]byte("test")))
		}()
	}
	wg.Wait()
}

func TestCacheRepository_Get(t *testing.T) {
	repo := &cacheRepository{cache: map[string]*cache{}}
	repo.Set("a", newCache([]byte("test")))
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			repo.Get("a")
		}()
	}
	wg.Wait()

	v, ok := repo.Get("a")
	if !ok {
		t.Errorf("Expected ok, but no")
	}
	if string(v.Bytes()) != "test" {
		t.Errorf("Expected test, got %v", string(v.Bytes()))
	}
}

func TestCacheRepository_Vacuum(t *testing.T) {
	repo := &cacheRepository{cache: map[string]*cache{}}
	repo.Set("a", newCache([]byte("test a")))
	repo.Set("b", newCache([]byte("test b")))
	repo.Set("c", newCache([]byte("test c")))

	tmp := timeutil.Now
	timeutil.Now = func() time.Time {
		return tmp().Add(cacheExpiredInterval * time.Second)
	}
	defer func() {
		timeutil.Now = tmp
	}()

	repo.Vacuum()

	if _, ok := repo.Get("a"); ok {
		t.Errorf("Expected a key deleted, but exists")
	}
	if _, ok := repo.Get("b"); ok {
		t.Errorf("Expected b key deleted, but exists")
	}
	if _, ok := repo.Get("c"); ok {
		t.Errorf("Expected c key deleted, but exists")
	}
}
