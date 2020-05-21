package url_shortener

import (
	"database/sql"
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

type PersistentUrlRepository struct {
	db sql.DB
}

func (c *CachedUrlRepository) Get(key string) (string, error) {
	val, err := c.client.Get(key).Result()
	if err != nil {
		log.Printf("Unable to find value for %s", key)
		return "", err
	}

	return val, nil
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

func (p *PersistentUrlRepository) Get(key string) (string, error) {
	query := "select long from url where key = :key"
	r, err := p.db.Query(query, sql.Named("key", key))
	if err != nil {
		panic(err)
	}
	defer r.Close()

	var value string
	err = r.Scan(&value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (p *PersistentUrlRepository) Set(key string, value string) {
	query := "insert into url (key, value) values (:key, :value)"
	_, err := p.db.Exec(query, sql.Named("key", key), sql.Named("value", value))
	if err != nil {
		panic(err)
	}
}
