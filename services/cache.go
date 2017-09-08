package services

import (
	"bytes"
	"errors"
	"time"

	"encoding/gob"

	"github.com/fpay/gopress"
	"github.com/garyburd/redigo/redis"
)

const (
	// CacheServiceName is the identity of cache service
	CacheServiceName = "cache"
)

type Index uint8

const (
	// User is the code of user index in redis
	User Index = iota + 1
	// Post is the code of post index in redis
	Post
	// Comment is the code of comment db in redis
	Comment
	// Tag is the code of user tag in redis
	Tag
	// Test is the code of test tag in redis
	Test
)

// CacheService type
type CacheService struct {
	// Uncomment this line if this service has dependence on other services in the container
	// c *gopress.Container
	cachePool *redis.Pool
}

// NewCacheService returns instance of cache service
func NewCacheService() *CacheService {
	p := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   1024,
		IdleTimeout: time.Duration(time.Second * 10),
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute*2 {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return &CacheService{
		cachePool: p,
	}
}

// ServiceName is used to implements gopress.Service
func (s *CacheService) ServiceName() string {
	return CacheServiceName
}

// RegisterContainer is used to implements gopress.Service
func (s *CacheService) RegisterContainer(c *gopress.Container) {
}

func (s *CacheService) encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func (s *CacheService) decode(serialData []byte, entity interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(serialData))
	if err := decoder.Decode(entity); err != nil {
		return err
	}
	return nil
}

// Get is that implement of redis GET operation
func (s *CacheService) Get(key string, entity interface{}, index Index) error {
	c := s.cachePool.Get()
	if c == nil {
		return errors.New("Can't get cache cenction")
	}
	defer c.Close()
	if _, err := c.Do("SELECT", index); err != nil {
		return err
	}
	reply, err := c.Do("GET", key)
	serialData, err := redis.Bytes(reply, err)
	if err != nil {
		return err
	}
	return s.decode(serialData, entity)
}

// Set is that implement of redis SET operation
func (s *CacheService) Set(key string, entity interface{}, index Index) error {
	c := s.cachePool.Get()
	if c == nil {
		return errors.New("Can't get cache cenction")
	}
	defer c.Close()
	if _, err := c.Do("SELECT", index); err != nil {
		return err
	}
	if serialData, err := s.encode(entity); err != nil {
		return err
	} else if _, err := c.Do("SET", key, serialData); err != nil {
		return err
	}
	return nil
}

// Del is that implement of redis  DEL operation
func (s *CacheService) Del(key string, index Index) error {
	c := s.cachePool.Get()
	if c == nil {
		return errors.New("Can't get cache cenction")
	}
	defer c.Close()
	if _, err := c.Do("SELECT", index); err != nil {
		return err
	}
	if _, err := c.Do("DEL", key); err != nil {
		return err
	}
	return nil
}

// SetExpireTime is that implement of redis HashSet action
func (s *CacheService) SetExpireTime(key string, interval int) error {
	c := s.cachePool.Get()
	if c == nil {
		return errors.New("Can't get cache cenction")
	}
	defer c.Close()
	if _, err := c.Do("EXPIRE", key, interval); err != nil {
		return err
	}
	return nil
}
