package main

import (
	"github.com/astaxie/beego"
	"github.com/barnettZQG/alertCenter/core"
	_ "github.com/barnettZQG/alertCenter/routers"
)

func main() {
	re := core.Relation{}
	beego.AddAPPStartHook(func() error {
		return re.Init(&core.LDAPServer{})
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
