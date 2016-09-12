package controllers

type UserController struct {
	BaseController
}

func (e *UserController) UserHome() {
	userName := e.GetString(":userName")
	e.Data["userName"] = userName
	e.TplName = "userHome.html"
}
