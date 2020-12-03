package ftconf

import (
	stdlog "log"
	"os"
)

const (
	Debug   = "debug"
	Test    = "test"
	Release = "release"
)

var (
	Mode = Debug
)

func init() {
	log := stdlog.New(os.Stdout, "[conf] ", stdlog.LstdFlags)

	m := os.Getenv("FTGO_MODE")
	if m == Debug ||
		m == Test ||
		m == Release {
		Mode = m
	}

	log.Printf("FTGO_MODE: %s    Mode: %s", m, Mode)
}
