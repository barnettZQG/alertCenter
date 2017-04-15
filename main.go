package main

import (
	"alertCenter/core/db"
	"alertCenter/core/notice"
	"alertCenter/core/user"
	_ "alertCenter/routers"
	_ "net/http/pprof"

	"alertCenter/core/service"

	"github.com/astaxie/beego"
)

func main() {
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
	beego.Run()
}
