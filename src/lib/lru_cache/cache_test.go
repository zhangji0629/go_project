package lru_cache_test

import (
	"lib/lru_cache"
	"strconv"
	"testing"
	"time"
)

func TestGet1(t *testing.T) {
	for i := 1; i <= 10; i++ {
		key := "name_" + strconv.FormatInt(int64(i), 10)
		value := i
		lru_cache.SetEx(key, value, time.Second*10)
	}

	time.Sleep(5 * time.Second)

	lru_cache.Del("name_3")
	lru_cache.Del("name_5")

	for i := 1; i <= 10; i++ {
		key := "name_" + strconv.FormatInt(int64(i), 10)
		t.Log(lru_cache.Get(key))
	}

	time.Sleep(time.Second * 6)

	t.Log("################################")
	for i := 1; i <= 10; i++ {
		key := "name_" + strconv.FormatInt(int64(i), 10)
		t.Log(lru_cache.Get(key))
	}
}

func TestGet2(t *testing.T) {
	lru_cache.SetEx("name", "zhangji", time.Second*10)
	t.Log(lru_cache.Get("name"))
}

func TestGet3(t *testing.T) {
	t.Log(
		lru_cache.GetWithFunc("10", 5*time.Second, func() (interface{}, error) {
			p := fn("10")
			return p, nil
		}))

	t.Log(lru_cache.Get("10"))
	time.Sleep(time.Second * 6)
	t.Log(lru_cache.Get("10"))
}

func fn(key string) string {
	time.Sleep(time.Millisecond * 50)
	return key + "____||||"
}

func BenchmarkGet1(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		if ret, _ := lru_cache.Get("10"); ret == nil {
			p := fn("10")
			lru_cache.SetEx("10", p, time.Minute)
		}
	}
}

func BenchmarkGet2(b *testing.B) {
	for i := 1; i <= b.N; i++ {
		lru_cache.GetWithFunc("10", time.Minute, func() (interface{}, error) {
			p := fn("10")
			return p, nil
		})
	}
}
