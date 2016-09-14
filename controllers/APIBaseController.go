package controllers

import (
	"alertCenter/core/db"
	"alertCenter/core/service"

	"alertCenter/util"

	"github.com/astaxie/beego"
)

type APIBaseController struct {
	beego.Controller
	Username string
}

func (this *APIBaseController) Prepare() {
	token := this.Ctx.Input.Header("token")
	user := this.Ctx.Input.Header("user")
	tokenService := &service.TokenService{
		Session: db.GetMongoSession(),
	}
	beego.Debug("check token:" + token + ",user:" + user)
	if ok := tokenService.CheckToken(token, user); !ok {
		this.Data["json"] = util.GetErrorJson("Security verification failed")
		this.ServeJSON()
	}
	this.Username = user
}
