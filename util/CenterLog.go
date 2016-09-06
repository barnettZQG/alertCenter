package util

import (
	"time"

	"github.com/astaxie/beego/logs"
)

var log logs.Logger = logs.NewConsole()

func Info(message string) {
	log.WriteMsg(time.Now(), message, logs.LevelInfo)
}

func Debug(message string) {
	log.WriteMsg(time.Now(), message, logs.LevelDebug)
}
func Error(message string) {
	log.WriteMsg(time.Now(), message, logs.LevelError)
}
