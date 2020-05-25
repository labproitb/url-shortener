package main

import (
	"database/sql"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"url-shortener/url"
)

var useCase url.UseCase
var cache url.CachedRepository
var storage url.PersistentRepository

func initializeCache() {
	r := redis.NewClient(&redis.Options{
		Addr:               os.Getenv("REDIS_ADDR"),
		Password:           "",
		DB:                 0,
	})
	cache = url.CachedRepository{Client:r}
}

func initializeStorage() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	storage = url.PersistentRepository{DB:db}
}

func main() {
	initializeCache()
	initializeStorage()

	useCase = url.UseCase{
		Cache:   &cache,
		Storage: &storage,
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/url", ShortenUrl)

	e.Logger.Fatal(e.Start(":4000"))
}

type ShortenUrlRequest struct {
	OriginalUrl string `json:"originalUrl"`
	ShortUrl 	string `json:"shortUrl"`
}

type ShortenUrlResponse struct {
	ShortUrl	string `json:"shortUrl"`
}

func ShortenUrl(c echo.Context) error {
	req := new(ShortenUrlRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	shortUrl := req.ShortUrl
	if shortUrl == "" {
		shortUrl = useCase.GenerateShortUrl()
	}
	useCase.Save(req.OriginalUrl, shortUrl)

	res := ShortenUrlResponse{ShortUrl:shortUrl}
	return c.JSON(http.StatusCreated, res)
}
