package controllers

import (
	"alertCenter/core/db"
	"alertCenter/core/service"
	"alertCenter/util"
)

type TokenAPIController struct {
	APIBaseController
}

func (e *TokenAPIController) GetAllToken() {

	user := e.Ctx.Input.Header("user")
	if user == "" {
		e.Data["json"] = util.GetErrorJson("please certification")
		e.ServeJSON()
	} else {
		service := &service.TokenService{
			Session: db.GetMongoSession(),
		}
		tokens := service.GetAllToken(user)
		e.Data["json"] = util.GetSuccessReJson(tokens)
		e.ServeJSON()
	}
}
