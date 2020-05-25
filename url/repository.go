package url

import (
	"database/sql"
	"github.com/go-redis/redis/v7"
	"log"
	"time"
)

type Repository interface {
	Get(key string) string
	Set(url string, key string) string
}

type CachedRepository struct {
	Client *redis.Client
}

type PersistentRepository struct {
	DB *sql.DB
}

func (c *CachedRepository) Get(key string) (string, error) {
	val, err := c.Client.Get(key).Result()
	if err != nil {
		log.Printf("Unable to find value for %s", key)
		return "", err
	}

	return val, nil
}

func (c *CachedRepository) Set(url string, key string) {
	d, err := time.ParseDuration("168h")
	if err != nil {
		d = time.Hour * 168
	}

	_, err = c.Client.Set(key, url, d).Result()
	if err != nil {
		panic(err)
	}
}

func (p *PersistentRepository) Get(key string) (string, error) {
	query := "select long from url where key = :key"
	r, err := p.DB.Query(query, sql.Named("key", key))
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

func (p *PersistentRepository) Set(key string, value string) {
	query := "insert into url (key, value) values (:key, :value)"
	_, err := p.DB.Exec(query, sql.Named("key", key), sql.Named("value", value))
	if err != nil {
		panic(err)
	}
}
