package controllers


type TeamController struct {
	BaseController
}

func (e *TeamController) GetTeams() {
	e.TplName = "teams.html"
}
