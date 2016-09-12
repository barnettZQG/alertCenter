package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
)

type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	beego.Debug("in /")
	code := c.GetString("code")
	beego.Debug("code is :", code)
	if code != "" {
		beego.Debug("code is:", code)
	} else {
		beego.Debug("do not have code.")
	}

	//c.Data["Website"] = "beego.me"
	c.TplName = "index.html"
}

type Transit struct {
	Response http.ResponseWriter
	Request  *http.Request
	Redirct  string
	code     int
}

func (c *MainController) Transit() {
	transit := c.GetSession("transit").(*Transit)
	if transit != nil {
		http.Redirect(transit.Response, transit.Request, transit.Redirct, http.StatusTemporaryRedirect)
	}
	redirct := beego.AppConfig.String("url")
	http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, redirct, http.StatusTemporaryRedirect)
}
