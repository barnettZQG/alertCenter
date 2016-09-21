package controllers

import (
	"alertCenter/core/service"

	"alertCenter/util"

	"github.com/astaxie/beego"
)

type APIBaseController struct {
	beego.Controller
	Username string
}

func (this *APIBaseController) Prepare() {
	// start := time.Now()
	// defer fmt.Println("check token time:", time.Now().Sub(start))

	token := this.Ctx.Input.Header("token")
	user := this.Ctx.Input.Header("user")
	tokenService := &service.TokenService{
		Session: nil,
	}
	beego.Debug("check token:" + token + ",user:" + user)
	if ok := tokenService.CheckToken(token, user); !ok {
		this.Data["json"] = util.GetErrorJson("Security verification failed")
		this.ServeJSON()
	}
	this.Username = user
}
