package main

import (
	"alertCenter/core"
	"alertCenter/core/notice"
	_ "alertCenter/routers"

	"github.com/astaxie/beego"
)

func main() {
	re := core.Relation{}
	beego.AddAPPStartHook(func() error {
		return re.Init()
	})
	beego.AddAPPStartHook(func() error {
		return notice.StartCenter()
	})
	beego.Run()
}
