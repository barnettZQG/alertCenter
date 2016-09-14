package controllers

import (
	"alertCenter/controllers/session"
	"net/http"

	"github.com/astaxie/beego"
	"alertCenter/util"
)

type MainController struct {
	BaseController
}


func (c *MainController) Logout(){
	sess, err := session.GetSession(c.Ctx)
	if err!=nil{
		beego.Error(err)
		return
	}
	defer sess.SessionRelease(c.Ctx.ResponseWriter)

	sess.Delete(session.SESSION_USERNAME)
	sess.Delete(session.SESSION_USER)

	c.Data["json"] = util.GetSuccessJson("Logout success.")
	c.ServeJSON()

	return
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
	sess, err := session.GetSession(c.Ctx)
	if err == nil {
		transit := sess.Get(session.SESSION_TRANSIT)
		if transit != nil {
			http.Redirect(transit.(*Transit).Response, transit.(*Transit).Request, transit.(*Transit).Redirct, http.StatusTemporaryRedirect)
		}
	}
	redirct := beego.AppConfig.String("url")
	http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, redirct, http.StatusTemporaryRedirect)
}
