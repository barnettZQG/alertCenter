package controllers

import (
	"alertCenter/controllers/session"
	"alertCenter/core/db"
	"alertCenter/core/service"
	"alertCenter/models"
	"alertCenter/util"

	"github.com/astaxie/beego"
)

type TokenAPIController struct {
	APIBaseController
}

func (e *TokenAPIController) GetAllToken() {
	sess, err := session.GetSession(e.Ctx)
	if err != nil {
		beego.Error("get session error:", err)
		e.Data["json"] = util.GetErrorJson("please certification")
		e.ServeJSON()
	}
	user := sess.Get(session.SESSION_USER)
	if user == nil {
		e.Data["json"] = util.GetErrorJson("please certification")
		e.ServeJSON()
	} else {
		service := &service.TokenService{
			Session: db.GetMongoSession(),
		}
		tokens := service.GetAllToken(user.(*models.User).Name)
		e.Data["json"] = util.GetSuccessReJson(tokens)
		e.ServeJSON()
	}
}
