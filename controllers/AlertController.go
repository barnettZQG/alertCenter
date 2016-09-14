package controllers

type AlertController struct {
	BaseController
}

func (e *AlertController) AlertList() {
	receiver := e.GetString("receiver")
	if len(receiver) != 0 {
		e.Data["receiver"] = receiver
		e.TplName = "alertList.html"
	} else {
		e.TplName = "alertListAll.html"
	}
}


