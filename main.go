package main

import (
	"alertCenter/core/db"
	"alertCenter/core/notice"
	"alertCenter/core/user"
	_ "alertCenter/routers"
	"log"
	"net/http"
	_ "net/http/pprof"

	"alertCenter/core/service"

	"github.com/astaxie/beego"
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
	//初始化检查全局配置
	beego.AddAPPStartHook(func() error {
		service := &service.GlobalConfigService{
			Session: db.GetMongoSession(),
		}
		if service.Session != nil {
			defer service.Session.Close()
		}
		return service.Init()
	})
	//beego.SetLogger("file", `{"filename":"log/test.log","level":10}`)
	beego.Info("mongo:", beego.AppConfig.String("mongoURI"))
	beego.Run()
}
