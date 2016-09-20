package controllers

type AlertController struct {
	BaseController
}

func (e *AlertController) AlertList() {
	e.TplName = "alertListAll.html"
}

func (e *AlertController) AlertsCurrent() {
	e.TplName = "alertList.html"
}
