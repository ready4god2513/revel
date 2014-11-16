package cache

import (
	"github.com/revel/revel"
	"net"
	_ "strings"
	"testing"
	"time"
)

// These tests require redis server running on localhost:6379 (the default)
const redisTestServer = "localhost:6379"
const incrementKey = "website-hits"
const incrementCount = 10

var newRedisCache = func(t *testing.T, defaultExpiration time.Duration) Cache {
	revel.Config = revel.NewEmptyConfig()

	c, err := net.Dial("tcp", redisTestServer)
	if err == nil {
		c.Write([]byte("flush_all\r\n"))
		c.Close()
		redisCache := NewRedisCache(redisTestServer, "", defaultExpiration)
		redisCache.Flush()
		return redisCache
	}
	t.Errorf("couldn't connect to redis on %s", redisTestServer)
	t.FailNow()
	panic("")
}

func increment(cache Cache, expected uint64, t *testing.T) {
	newValue, err := cache.Increment(incrementKey, incrementCount)

	if err != nil {
		t.Errorf("Error incrementing redis key %s: %s", incrementKey, err)
	}

	if newValue != expected {
		t.Errorf("Expected %d, was %d", newValue, expected)
	}
}

func TestRedisCache_TypicalGetSet(t *testing.T) {
	typicalGetSet(t, newRedisCache)
}

func TestRedisCache_IncrementOnUndefined(t *testing.T) {
	cache := newRedisCache(t, time.Minute)
	increment(cache, incrementCount, t)
}

func TestRedisCache_Increment(t *testing.T) {
	cache := newRedisCache(t, time.Minute)
	cache.Set(incrementKey, incrementCount, time.Hour)
	increment(cache, incrementCount*2, t)
}

func TestRedisCache_DecrementOnUndefined(t *testing.T) {
	cache := newRedisCache(t, time.Minute)
	increment(cache, incrementCount, t)
}

func TestRedisCache_Decrement(t *testing.T) {
	cache := newRedisCache(t, time.Minute)
	cache.Set(incrementKey, incrementCount, time.Hour)
	increment(cache, incrementCount*2, t)
}

func TestRedisCache_Expiration(t *testing.T) {
	expiration(t, newRedisCache)
}

func TestRedisCache_Replace(t *testing.T) {
	testReplace(t, newRedisCache)
}

func TestRedisCache_Add(t *testing.T) {
	testAdd(t, newRedisCache)
}

func TestRedisCache_GetMulti(t *testing.T) {
	testGetMulti(t, newRedisCache)
}
