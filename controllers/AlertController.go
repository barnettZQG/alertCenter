package controllers

import "github.com/astaxie/beego"

type AlertController struct {
	beego.Controller
}

func (e *AlertController) AlertList() {
	receiver := e.GetString("receiver")
	if len(receiver) != 0 {
		e.TplName = "alertList.html"
	} else {
		e.TplName = "alertListAll.html"
	}
}
