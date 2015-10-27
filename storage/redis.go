package storage

import (
	"net/url"

	"github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/gopkg.in/redis.v3"
)

type RedisClient redis.Client

const key = "time"

func NewRedisClient(urlS string, initialTime int64) (*RedisClient, error) {
	redisURL, err := url.Parse(urlS)
	if err != nil {
		return nil, err
	}
	redisPassword, _ := redisURL.User.Password()
	redisOptions := redis.Options{
		Addr:     redisURL.Host,
		Password: redisPassword,
		DB:       0,
	}
	rc := RedisClient(*redis.NewClient(&redisOptions))
	_, err = rc.Ping().Result()
	if err != nil {
		return nil, err
	}
	if initialTime != 0 {
		return &rc, rc.SetTime(initialTime)
	}
	return &rc, err
}

func (rc *RedisClient) SetTime(time int64) error {
	return rc.Set(key, time, 0).Err()
}

func (rc *RedisClient) GetTime() (int64, error) {
	return rc.Get(key).Int64()
}
