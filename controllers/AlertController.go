package controllers

import (
	"alertCenter/core/db"
	"alertCenter/core/service"
)

type AlertController struct {
	BaseController
}

func (e *AlertController) AlertList() {
	e.TplName = "alertListAll.html"
}

func (e *AlertController) AlertsCurrent() {
	e.TplName = "alertList.html"
}

func (e *AlertController) HistoryList() {
	mark := e.GetString(":mark")
	if mark == "" {
		e.Abort("404")
	}
	alertService := &service.AlertService{
		Session: db.GetMongoSession(),
	}
	if alertService.Session != nil {
		defer alertService.Session.Close()
	}
	alert, _ := alertService.GetAlertByMark(mark)
	if alert == nil {
		e.Abort("404")
	} else {
		e.Data["alert"] = alert
		e.Data["mark"] = mark
		e.Data["description"] = alert.Annotations.LabelSet["description"]
		e.Data["alertName"] = alert.Labels.LabelSet["alertname"]
		e.Data["startsAt"] = alert.StartsAt.Format("2006-01-02T15:04:05")
		e.Data["endsAt"] = alert.EndsAt.Format("2006-01-02T15:04:05")
		labels := make(map[string]string, 0)
		for k := range alert.Labels.LabelSet {
			labels[string(k)] = string(alert.Labels.LabelSet[k])
		}
		e.Data["labels"] = labels
		e.TplName = "historyList.html"
	}
}
