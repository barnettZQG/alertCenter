package main

import (
	"alertCenter/core"
	_ "alertCenter/routers"

	"github.com/astaxie/beego"
)

func main() {
	re := core.Relation{}
	beego.AddAPPStartHook(func() error {
		return re.Init()
	})
	// we := core.WeAlertSend{}
	// beego.AddAPPStartHook(func() error {
	// 	return we.StartWork()
	// })
	// mail := core.MailAlertSend{}
	// beego.AddAPPStartHook(func() error {
	// 	return mail.StartWork()
	// })
	beego.Run()
	//mail.StopWork()
	//we.StopWork()
}
