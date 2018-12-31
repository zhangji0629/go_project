package lru_cache_test

import (
	"lru_cache"
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
