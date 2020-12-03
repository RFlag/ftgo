package ftsync

import (
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/yesuu/redsync"
)

const prefix = "redlock:"

var redsync_ *redsync.Redsync

func init() {
	addr := os.Getenv("FTGO_REDIS")
	auth := os.Getenv("FTGO_REDIS_AUTH")
	log.Printf("FTGO_REDIS: %s    FTGO_REDIS_AUTH: %s", addr, auth)

	redsync_ = redsync.New([]redsync.Pool{
		&redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", addr)
				if err != nil {
					return nil, err
				}
				if _, err := c.Do("AUTH", auth); err != nil {
					c.Close()
					return nil, err
				}
				return c, nil
			},
		},
	})
}

// 阻塞锁

func Lock(name string, options ...redsync.Option) (*redsync.Mutex, error) {
	m := redsync_.NewMutex(prefix+name, options...)
	err := m.Lock()
	return m, err
}

// 非阻塞锁
func OptimisticLock(name string, options ...redsync.Option) (*redsync.Mutex, error) {
	m := redsync_.NewMutex(prefix+name, options...)
	err := m.OptimisticLock()
	return m, err
}
