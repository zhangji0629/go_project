package lru_cache

import (
	"errors"
	"lib/singleflight"
	"time"

	"github.com/cespare/xxhash"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sirupsen/logrus"
)

var (
	dbs    map[byte]*lru.Cache
	Nil    = errors.New("cache:nil")
	gGroup = new(singleflight.Group)
)

type item struct {
	v interface{}
	d int64
}

func init() {
	if dbs == nil {
		dbs = make(map[byte]*lru.Cache)
		for i := 0; i < 256; i++ {
			dbs[byte(i)], _ = lru.New(20000)
		}
	}
	go flush()
}

func flush() {
	i := 0
	for range time.NewTicker(time.Minute).C {
		shard := byte(i % 256)
		db, _ := dbs[shard]
		keys := db.Keys()
		i++
		for _, key := range keys {
			if keyS, ok := key.(string); ok {
				getByDb(keyS, db)
			}
		}
		logrus.Infof("FlushKey||time=%s||shard=%d||old=%d||new=%d\n", time.Since(now), shard, len(keys), db.Len())
	}
}

func SetEx(key string, value interface{}, ex time.Duration) error {
	return setEx(key, value, ex)
}

func Set(key string, value interface{}) error {
	return setEx(key, value, time.Hour*24)
}

func Get(key string) (interface{}, error) {
	db, _ := dbs[byte(hash(key))]
	return getByDb(key, db)
}

func FreqCall(key string, ex time.Duration, fn func()) {
	if value, _ := Get(key); value == nil {
		setEx(key, 1, ex)
		fn()
	}
}

func GetWithFunc(key string, ex time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	if tmp, _ := Get(key); tmp != nil {
		return tmp, nil
	}
	ret, err := gGroup.Do(key, fn)
	if err != nil {
		return nil, err
	}
	setEx(key, ret, ex)
	return ret, nil
}

func Del(key string) error {
	db, _ := dbs[hash(key)]
	delByDb(key, db)
	return nil
}

func setEx(key string, value interface{}, ex time.Duration) error {
	db, _ := dbs[byte(hash(key))]
	db.Add(key, &item{
		v: value,
		d: time.Now().Add(ex).UnixNano(),
	})
	return nil
}

func getByDb(key string, db *lru.Cache) (interface{}, error) {
	if value, ok := db.Get(key); ok {
		if item, ok := value.(*item); ok {
			if time.Now().UnixNano() < item.d {
				return item.v, nil
			}
			delByDb(key, db)
		}
	}
	return nil, Nil
}

func delByDb(key string, db *lru.Cache) error {
	db.Remove(key)
	return nil
}

func hash(key string) byte {
	return byte(xxhash.Sum64String(key))
}
