package ftredis

import (
	stdlog "log"
	"os"

	"ftgo/safeclose"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var Client *redis.Client

func init() {
	log := stdlog.New(os.Stdout, "[redis] ", stdlog.LstdFlags)

	addr := os.Getenv("FTGO_REDIS")
	auth := os.Getenv("FTGO_REDIS_AUTH")
	log.Printf("FTGO_REDIS: %s    FTGO_REDIS_AUTH: %s", addr, auth)

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: auth,
	})

	pong, err := Client.Ping().Result()
	if err != nil {
		safeclose.Cancel()
		log.Print(errors.Wrap(err, "pong="+pong))
		safeclose.Wait()
		os.Exit(1)
		return
	}
	log.Print(pong)

	safeclose.Defer(func() {
		err := Client.Close()
		if err != nil {
			log.Print(errors.Wrap(err, "关闭 redis 失败"))
			return
		}
		log.Print("redis 安全关闭")
	})
}
