package url

import "math/rand"

type UseCase struct {
	Cache CachedUrlRepository
	Storage PersistentUrlRepository
}

func (u *UseCase) generateRandomString(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, length)
		for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func (u *UseCase) GenerateShortUrl() string {
	var s string
	for {
		s = u.generateRandomString(5)
		_, err := u.Storage.Get(s)
		if err != nil {
			break
		}
	}
	return s
}

func (u *UseCase) Save(url string, shortUrl string) {
	u.Storage.Set(url, shortUrl)
}

func (u *UseCase) Get(shortUrl string) (string, error) {
	url, err := u.Cache.Get(shortUrl)
	if err == nil {
		return url, nil
	}

	url, err = u.Storage.Get(shortUrl)
	if err != nil {
		return "", err
	}

	u.Cache.Set(url, shortUrl)
	return url, nil
}

