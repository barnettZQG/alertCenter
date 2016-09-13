package controllers

import (
	"alertCenter/controllers/session"
	"alertCenter/core/db"
	"alertCenter/core/service"
	"alertCenter/models"
	"alertCenter/util"

	"github.com/astaxie/beego"
)

type TokenController struct {
	BaseController
}

func (e *TokenController) AddToken() {
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
		u := user.(*models.User)
		projectName := e.GetString("projectName")
		if projectName == "" {
			e.Data["json"] = util.GetErrorJson("projectName can't be empty")
			e.ServeJSON()
		} else {
			service := &service.TokenService{
				Session: db.GetMongoSession(),
			}
			token := service.CreateToken(projectName, u.Name)
			e.Data["json"] = util.GetSuccessReJson(token)
			e.ServeJSON()
		}
	}
}
