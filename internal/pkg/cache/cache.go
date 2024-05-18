package cache

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

type CacheItem struct {
	Key     string
	Value   []byte
	ExpTime int32
}

type Config struct {
	Host                string `yaml:"host"`
	Port                string `yaml:"port"`
	OrderExpirationTime int    `yaml:"order_expiration_time"`
}

type Cache interface {
	Set(item *CacheItem) error
	Get(key string) (*CacheItem, error)
}

type memCached struct {
	client *memcache.Client
	cfg    *Config
}

func (m *memCached) Set(item *CacheItem) error {
	return m.client.Set(&memcache.Item{
		Key:        item.Key,
		Value:      item.Value,
		Expiration: item.ExpTime,
	})
}

func (m *memCached) Get(key string) (*CacheItem, error) {
	mItem, err := m.client.Get(key)

	if err != nil {
		return nil, err
	}

	return &CacheItem{
		Key:     mItem.Key,
		Value:   mItem.Value,
		ExpTime: mItem.Expiration,
	}, nil
}

func NewCache(cfg *Config) Cache {
	return &memCached{
		client: memcache.New(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
		cfg:    cfg,
	}
}