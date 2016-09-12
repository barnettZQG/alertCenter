package controllers

import (
	"alertCenter/core/user"
	"alertCenter/models"
)

type UserController struct {
	BaseController
}

func (e *UserController) UserHome() {
	userName := e.GetString(":userName")
	relaction := user.Relation{}
	user := relaction.GetUserByName(userName)
	self := e.GetSession("user")
	if user == nil {
		e.Abort("404")
	} else {
		e.Data["user"] = user
		if self != nil && user.Name == self.(*models.User).Name {
			e.Data["self"] = true
		}
		e.TplName = "userHome.html"
	}
}
