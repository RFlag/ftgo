package ftsql

import (
	stdlog "log"
	"os"

	"ftgo/safeclose"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var DB *sqlx.DB

func init() {
	log := stdlog.New(os.Stdout, "[sql] ", stdlog.LstdFlags)

	url := os.Getenv("FTGO_SQL")
	log.Print("FTGO_SQL: " + url)

	var err error
	DB, err = sqlx.Connect("mysql", url)
	if err != nil {
		safeclose.Cancel()
		log.Print(errors.Wrap(err, "sql 连接数据库失败"))
		safeclose.Wait()
		os.Exit(1)
		return
	}

	safeclose.Defer(func() {
		err := DB.Close()
		if err != nil {
			log.Print(errors.Wrap(err, "关闭 sql 失败"))
			return
		}
		log.Print("sql 安全关闭")
	})
}
