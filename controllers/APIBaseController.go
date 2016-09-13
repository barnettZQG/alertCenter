package controllers

import (
	"alertCenter/core/db"
	"alertCenter/core/service"

	"alertCenter/util"

	"github.com/astaxie/beego"
)

type APIBaseController struct {
	beego.Controller
}

func (this *APIBaseController) Prepare() {
	token := this.GetString("token")
	user := this.GetString("user")
	tokenService := &service.TokenService{
		Session: db.GetMongoSession(),
	}
	if ok := tokenService.CheckToken(token, user); !ok {
		this.Data["json"] = util.GetErrorJson("Security verification failed")
		this.ServeJSON()
	}
}
