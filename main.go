package main

import (
	"alertCenter/core/notice"
	"alertCenter/core/user"
	_ "alertCenter/routers"

	"github.com/astaxie/beego"
)

func init() {

}

func main() {
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
