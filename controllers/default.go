package controllers

import "github.com/astaxie/beego"

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	beego.Debug("in /")
	code := c.GetString("code")
	beego.Debug("code is :",code)
	if code !=""{
		beego.Debug("code is:",code)
	}else{
		beego.Debug("do not have code.")
	}

	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
