package ftmongo

import (
	stdlog "log"
	"os"
	"time"

	"ftgo/safeclose"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

const dbname = "ranbb"

var DB *mgo.Database

func init() {
	log := stdlog.New(os.Stdout, "[mgo] ", stdlog.LstdFlags)

	url := os.Getenv("FTGO_MONGO")
	log.Print("FTGO_MONGO: " + url)

	session, err := mgo.Dial(url)
	if err != nil {
		safeclose.Cancel()
		log.Print(errors.Wrap(err, "mgo 连接数据库失败"))
		safeclose.Wait()
		os.Exit(1)
		return
	}
	DB = session.DB(dbname)

	// mgo自动断线重连
	go func() {
		t := time.Tick(32 * time.Second)
		var closeErr error
		for {
			<-t
			closeErr = session.Ping()
			if closeErr != nil {
				log.Print(errors.Wrap(closeErr, "mgo ping不通"))
				session.Refresh()
				closeErr = session.Ping()
				if closeErr != nil {
					log.Print(errors.Wrap(closeErr, "mgo refresh失败"))
				} else {
					log.Print("mgo refresh成功")
				}
			}
		}
	}()

	safeclose.Defer(func() {
		session.Close()
		log.Print("mgo 安全关闭")
	})
}
