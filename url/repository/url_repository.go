package repository

import (
	"github.com/go-redis/redis/v7"
	"log"
	"time"
)

type UrlRepository interface {
	Get(key string) string
	Set(url string, key string) string
}

type CachedUrlRepository struct {
	client redis.Client
}

func (c *CachedUrlRepository) Get(key string) string {
	val, err := c.client.Get(key).Result()
	if err != nil {
		log.Printf("Unable to find value for %s", key)
	}

	return val
}

func (c *CachedUrlRepository) Set(url string, key string) {
	d, err := time.ParseDuration("168h")
	if err != nil {
		d = time.Hour * 168
	}

	_, err = c.client.Set(key, url, d).Result()
	if err != nil {
		panic(err)
	}
}
