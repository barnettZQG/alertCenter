package main

import (
	"alertCenter/core/notice"
	"alertCenter/core/user"
	_ "alertCenter/routers"
	_ "net/http/pprof"
	"github.com/astaxie/beego"
	"log"
	"net/http"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	re := user.Relation{}
	beego.AddAPPStartHook(func() error {
		return re.Init()
	})
	beego.AddAPPStartHook(func() error {
		return notice.StartCenter()
	})
	beego.Info("mongo:", beego.AppConfig.String("mongoURI"))
	beego.Run()
}
