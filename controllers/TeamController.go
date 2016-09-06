package controllers

import "github.com/astaxie/beego"

type TeamController struct {
	beego.Controller
}

func (e *TeamController) GetTeams() {
	e.TplName = "teams.html"
}
